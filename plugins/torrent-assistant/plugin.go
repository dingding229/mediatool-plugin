// Package assistant suggests by default and mutates only verified owned tasks through Host.
package assistant

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/dingding229/mediatool-plugin/plugins/internal/acceptance"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/downloads"
	"github.com/dingding229/mediatool-plugin/sdk/host/notifications"
	"math"
	"sync"
	"time"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Config struct {
	Active       bool               `json:"active"`
	Write        bool               `json:"write"`
	Recheck      bool               `json:"recheck"`
	Tag          string             `json:"tag"`
	Category     string             `json:"category"`
	DownloadKiB  acceptance.Numeric `json:"downloadKiB"`
	UploadKiB    acceptance.Numeric `json:"uploadKiB"`
	Ratio        acceptance.Numeric `json:"ratio"`
	SeedHours    acceptance.Numeric `json:"seedHours"`
	Batch        acceptance.Numeric `json:"batch"`
	StallMinutes acceptance.Numeric `json:"stallMinutes"`
}
type Sample struct {
	ID        string  `json:"id"`
	Hash      string  `json:"hash"`
	Progress  float64 `json:"progress"`
	Since     int64   `json:"since"`
	State     string  `json:"state"`
	Applying  bool    `json:"applying"`
	Advised   bool    `json:"advised"`
	Rechecked bool    `json:"rechecked"`
}
type State struct {
	Page   int      `json:"page"`
	Items  []Sample `json:"items"`
	Status string   `json:"status"`
}
type Plugin struct {
	mu      sync.Mutex
	host    *host.Client
	config  Config
	running bool
	now     func() time.Time
}

func New(s sdk.Services) sdk.Plugin {
	var _ sdk.ActionProvider = (*Plugin)(nil)
	return &Plugin{host: s.Host, now: time.Now}
}
func (p *Plugin) InitPlugin(_ context.Context, c sdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v := Config{Batch: 5, StallMinutes: 30}
	if contracts.Decode(c.Value, &v) != nil || len(v.Tag) > 64 || len(v.Category) > 64 {
		return contracts.NewError(contracts.InvalidParams)
	}
	for _, bound := range [][2]acceptance.Numeric{{v.DownloadKiB, 102400}, {v.UploadKiB, 102400}, {v.Ratio, 1000}, {v.SeedHours, 87600}, {v.Batch, 10}, {v.StallMinutes, 1440}} {
		n, max := bound[0], bound[1]
		if n < 0 || n > max || math.IsNaN(float64(n)) || math.IsInf(float64(n), 0) {
			return contracts.NewError(contracts.InvalidParams)
		}
	}
	if v.Batch < 1 || math.Trunc(float64(v.Batch)) != float64(v.Batch) || v.StallMinutes < 5 {
		return contracts.NewError(contracts.InvalidParams)
	}
	p.config = v
	p.running = true
	return nil
}
func (p *Plugin) GetState(context.Context) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running, nil
}
func (p *Plugin) StopService(context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.running = false
	return nil
}
func scopes() []string {
	return []string{"plugin.data.read", "plugin.data.write", "downloads.read", "downloads.control", "notifications.send"}
}
func (p *Plugin) GetService(context.Context) ([]sdk.Service, error) {
	return []sdk.Service{{ID: "assist", Name: "预览或执行下载辅助规则", Trigger: "interval", Interval: 5 * time.Minute, ScheduleSetting: "active", RequiredScopes: scopes(), Handler: p.run}}, nil
}
func positive(v acceptance.Numeric, mult int64) *int64 {
	if v == 0 {
		return nil
	}
	n := int64(float64(v) * float64(mult))
	return &n
}
func (p *Plugin) run(ctx context.Context, _ json.RawMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	s, rev, e := acceptance.Load[State](ctx, p.host, "assistant-state")
	if e != nil {
		return e
	}
	for _, v := range s.Items {
		if v.Applying {
			return acceptance.RequireAction(ctx, p.host, "assist-unknown", "assistant_outcome_unknown")
		}
	}
	if s.Page < 1 || s.Page > 1000 {
		s.Page = 1
	}
	c := downloads.New(p.host)
	page, e := c.List(ctx, b.PageRequest{Page: s.Page, Size: int(p.config.Batch)})
	if e != nil {
		return e
	}
	current := map[string]bool{}
	old := map[string]Sample{}
	for _, v := range s.Items {
		old[v.ID] = v
	}
	next := []Sample{}
	save := func() error {
		v, err := p.host.SaveData(ctx, "assistant-state", s, rev)
		if err == nil {
			rev = v.Revision
		}
		return err
	}
	for _, base := range page.Items {
		if !base.Owned || base.Hash == "" || base.SeedingState == "removed" {
			continue
		}
		t, err := c.Inspect(ctx, b.ID{ID: base.Hash})
		if err != nil {
			return err
		}
		if !t.Owned || t.ID != base.ID || t.Progress == nil {
			continue
		}
		v, ok := old[t.ID]
		if !ok || v.Hash != t.Hash || v.Progress != *t.Progress {
			v = Sample{ID: t.ID, Hash: t.Hash, Progress: *t.Progress, Since: p.now().Unix()}
		}
		problem := t.State == "error" || t.State == "missingFiles" || t.State == "stalledDL" && p.now().Unix()-v.Since >= int64(p.config.StallMinutes*60)
		if problem && !v.Advised {
			r, err := notifications.New(p.host).Send(ctx, b.Notice{Title: "下载器任务需关注", Text: "任务 " + t.ID + " 出错或进度静默；默认不自动重启。"}, b.Operation{Key: fmt.Sprintf("advice:%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", t.ID, t.Hash, v.Since)))), Confirm: true})
			if err != nil || r.State != b.NotificationPending {
				return contracts.NewError(contracts.ResultUnknown)
			}
			v.Advised = true
		}
		if !problem && t.State != "stalledDL" {
			v.Advised = false
			v.Rechecked = false
		}
		v.State = "preview"
		old[v.ID] = v
		current[v.ID] = true
	}
	for _, v := range old {
		next = append(next, v)
	}
	if len(next) > 200 {
		return acceptance.RequireAction(ctx, p.host, "assist-capacity", "assistant_ledger_full")
	}
	s.Items = next
	if page.HasMore {
		s.Page++
	} else {
		s.Page = 1
	}
	if e = save(); e != nil {
		return e
	}
	for i := range s.Items {
		v := &s.Items[i]
		if !current[v.ID] {
			continue
		}
		cfg := b.DownloaderConfig{ID: v.ID, Tag: p.config.Tag, Category: p.config.Category, DownloadLimit: positive(p.config.DownloadKiB, 1024), UploadLimit: positive(p.config.UploadKiB, 1024), SeedingMinutes: positive(p.config.SeedHours, 60)}
		if p.config.Ratio > 0 {
			ratio := float64(p.config.Ratio)
			cfg.RatioLimit = &ratio
		}
		if !p.config.Write {
			continue
		}
		configure := cfg.Tag != "" || cfg.Category != "" || cfg.DownloadLimit != nil || cfg.UploadLimit != nil || cfg.RatioLimit != nil || cfg.SeedingMinutes != nil
		raw, _ := json.Marshal(cfg)
		v.Applying = true
		if e = save(); e != nil {
			return e
		}
		if configure {
			r, err := c.Configure(ctx, cfg, b.Operation{Key: fmt.Sprintf("assist:%x", sha256.Sum256(raw)), Confirm: true})
			if err != nil || r.State != "configured" {
				return acceptance.RequireAction(ctx, p.host, "assist-unknown", "assistant_outcome_unknown")
			}
		}
		if p.config.Recheck && v.Advised && !v.Rechecked {
			r, err := c.Recheck(ctx, b.ID{ID: v.ID}, b.Operation{Key: fmt.Sprintf("recheck:%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", v.ID, v.Hash, v.Since)))), Confirm: true})
			if err != nil || r.State != "requested" {
				return acceptance.RequireAction(ctx, p.host, "assist-unknown", "assistant_outcome_unknown")
			}
			v.State = "recheck_requested"
			v.Rechecked = true
		} else {
			v.State = "configured"
		}
		v.Applying = false
		if e = save(); e != nil {
			return e
		}
	}
	s.Status = "ok"
	return save()
}
func (p *Plugin) GetActions(context.Context) ([]sdk.Action, error) {
	return []sdk.Action{{ID: "assist", Name: "危险：按已确认配置运行（默认仅预览）", RequiredScopes: scopes(), Handler: func(ctx context.Context, q sdk.ActionContext) (sdk.ActionContext, bool, error) {
		e := p.run(ctx, nil)
		return q, true, e
	}}}, nil
}
func (p *Plugin) GetCommand(context.Context) ([]sdk.Command, error) {
	return []sdk.Command{{Cmd: "downloader-status", Desc: "只读下载器助手状态", Category: "下载器", RequiredScopes: []string{"plugin.data.read"}, Handler: func(ctx context.Context, _ json.RawMessage) (json.RawMessage, error) {
		s, _, e := acceptance.Load[State](ctx, p.host, "assistant-state")
		raw, err := json.Marshal(s)
		if e == nil {
			e = err
		}
		return raw, e
	}}}, nil
}
func (p *Plugin) GetAPI(context.Context) ([]sdk.API, error) {
	return []sdk.API{{Path: "/status", Methods: []string{"GET"}, Auth: "admin", Summary: "只读下载器助手状态", RequiredScopes: []string{"plugin.data.read"}, Handler: func(ctx context.Context, _ sdk.APIRequest) (sdk.APIResponse, error) {
		s, _, e := acceptance.Load[State](ctx, p.host, "assistant-state")
		raw, _ := json.Marshal(s)
		return sdk.APIResponse{Status: 200, Body: raw}, e
	}}}, nil
}
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) {
	return sdk.Form{Defaults: json.RawMessage(`{"active":false,"write":false,"recheck":false,"batch":"5","stallMinutes":"30","downloadKiB":"0","uploadKiB":"0","ratio":"0","seedHours":"0"}`), Content: []contracts.Component{acceptance.Field("VSwitch", "active", "启用辅助检查"), acceptance.Field("VSwitch", "write", "调整下载设置", "危险：允许设置标签/分类/限制"), acceptance.Field("VSwitch", "recheck", "重新校验", "额外危险：允许错误/卡住任务重新校验（非重启进程）"), acceptance.Field("VTextField", "tag", "附加标签", "附加标签（不移除 owner 标签）"), acceptance.Field("VTextField", "category", "目标分类", "已有分类（自动管理路径模式拒绝）"), acceptance.Number("downloadKiB", "下载限速（KiB/s）", 0, 102400, 1, false, "下载 KiB/s，0=不修改"), acceptance.Number("uploadKiB", "上传限速（KiB/s）", 0, 102400, 1, false, "上传 KiB/s，0=不修改"), acceptance.Number("ratio", "分享率上限", 0, 1000, 0.1, false, "做种比例限制，0=不修改；HR 必须明确 false"), acceptance.Number("seedHours", "做种时长上限（小时）", 0, 87600, 0.1, false, "做种小时限制，0=不修改"), acceptance.Number("batch", "每轮任务上限", 1, 10, 1, true, "每轮任务数（1–10）"), acceptance.Number("stallMinutes", "停滞阈值（分钟）", 5, 1440, 1, false, "进度静默分钟（5–1440）")}}, nil
}
func (p *Plugin) GetPage(ctx context.Context, _ sdk.PageRequest) (sdk.Page, error) {
	s, _, e := acceptance.Load[State](ctx, p.host, "assistant-state")
	nodes := []contracts.Component{{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(fmt.Sprintf("状态 %s；任务 %d；显示前20项", s.Status, len(s.Items)))}}}
	for i, v := range s.Items {
		if i >= 20 {
			break
		}
		raw, _ := json.Marshal(v)
		nodes = append(nodes, contracts.Component{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(string(raw))}})
	}
	return sdk.Page{Content: nodes}, e
}

func (p *Plugin) GetDashboard(ctx context.Context, q sdk.DashboardRequest) (sdk.Dashboard, error) {
	return p.GetPage(ctx, q)
}
func (p *Plugin) GetDashboardMeta(context.Context) ([]sdk.DashboardMeta, error) {
	return []sdk.DashboardMeta{{Key: "overview", Title: "下载器助手（默认只提示）", Refresh: 60, Border: true, Cols: 6}}, nil
}
