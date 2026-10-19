// Package seeding previews owned, imported, known-HR-safe torrents; execution is opt-in.
package seeding

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sync"

	"github.com/dingding229/mediatool-plugin/plugins/internal/acceptance"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/downloads"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Config struct {
	Active          bool               `json:"active"`
	Execute         bool               `json:"execute"`
	DeleteFiles     bool               `json:"deleteFiles"`
	MinRatio        acceptance.Numeric `json:"minRatio"`
	MinSeedingHours acceptance.Numeric `json:"minSeedingHours"`
	BatchSize       acceptance.Numeric `json:"batchSize"`
	Category        string             `json:"category"`
	Tag             string             `json:"tag"`
	SiteID          string             `json:"siteId"`
}
type Item struct {
	TaskID string `json:"taskId"`
	Hash   string `json:"hash"`
	State  string `json:"state"`
}
type State struct {
	Page  int    `json:"page"`
	Items []Item `json:"items"`
}
type Plugin struct {
	mu      sync.Mutex
	host    *host.Client
	config  Config
	running bool
}

func New(s sdk.Services) sdk.Plugin { return &Plugin{host: s.Host} }
func (p *Plugin) InitPlugin(_ context.Context, c sdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v := Config{MinRatio: 2, MinSeedingHours: 168, BatchSize: 10}
	if contracts.Decode(c.Value, &v) != nil || (math.IsNaN(float64(v.MinSeedingHours)) || math.IsInf(float64(v.MinSeedingHours), 0) || math.IsNaN(float64(v.BatchSize)) || math.IsInf(float64(v.BatchSize), 0) || math.IsNaN(float64(v.MinRatio))) || math.IsInf(float64(v.MinRatio), 0) || v.MinRatio < 0 || v.MinRatio > 10000 || v.MinSeedingHours < 1 || v.MinSeedingHours > 87600 || v.BatchSize < 1 || v.BatchSize > 20 || math.Trunc(float64(v.BatchSize)) != float64(v.BatchSize) || len(v.Tag) > 200 || len(v.Category) > 200 || len(v.SiteID) > 200 {
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
func (p *Plugin) GetService(context.Context) ([]sdk.Service, error) {
	return []sdk.Service{{ID: "cleanup", Name: "预览或执行做种清理", Trigger: "cron", Cron: "*/10 * * * *", ScheduleSetting: "active", RequiredScopes: []string{"plugin.data.read", "plugin.data.write", "downloads.read", "downloads.delete"}, Handler: p.run}}, nil
}
func (p *Plugin) eligible(t b.Task) bool {
	c := p.config
	return t.Owned && t.ID != "" && t.Complete && t.Seeding && t.CleanupEligible && t.Ratio >= float64(c.MinRatio) && t.SeedingSeconds >= int64(c.MinSeedingHours*3600) && (c.Category == "" || c.Category == t.Category) && (c.Tag == "" || slices.Contains(t.Tags, c.Tag)) && (c.SiteID == "" || c.SiteID == t.SiteID)
}
func (p *Plugin) run(ctx context.Context, _ json.RawMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	s, rev, e := acceptance.Load[State](ctx, p.host, "cleanup-state")
	if e != nil {
		return e
	}
	if s.Page < 1 || s.Page > 1000 {
		s.Page = 1
	}
	// Intent is durable before any deletion. A crash/unknown needs an operator, not replay.
	for _, v := range s.Items {
		if v.State == "executing" || v.State == "unknown" {
			return acceptance.RequireAction(ctx, p.host, "cleanup-unknown", "cleanup_outcome_unknown")
		}
	}
	client := downloads.New(p.host)
	page, e := client.List(ctx, b.PageRequest{Page: s.Page, Size: 200})
	if e != nil {
		return e
	}
	s.Items = []Item{}
	for _, candidate := range page.Items {
		if len(s.Items) >= int(p.config.BatchSize) {
			break
		}
		if !candidate.Owned || candidate.Hash == "" || candidate.SeedingState == "removed" || candidate.CleanupState != "" && candidate.CleanupState != "disabled" {
			continue
		}
		t, err := client.Inspect(ctx, b.ID{ID: candidate.Hash})
		if err != nil {
			return err
		}
		if !p.eligible(t) || t.ID != candidate.ID {
			continue
		}
		item := Item{TaskID: t.ID, Hash: t.Hash, State: "preview"}
		s.Items = append(s.Items, item)
	}
	if page.HasMore {
		s.Page++
	} else {
		s.Page = 1
	}
	save := func() error {
		v, err := p.host.SaveData(ctx, "cleanup-state", s, rev)
		if err == nil {
			rev = v.Revision
		}
		return err
	}
	if e = save(); e != nil {
		return e
	}
	if !p.config.Execute {
		return nil
	}
	for i := range s.Items {
		s.Items[i].State = "executing"
		if e = save(); e != nil {
			return e
		}
		item := s.Items[i]
		key := fmt.Sprintf("seed:%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%t", item.TaskID, item.Hash, p.config.DeleteFiles))))
		r, err := client.DeleteOwned(ctx, b.DeleteRequest{ID: item.TaskID, DeleteFiles: p.config.DeleteFiles}, b.Operation{Key: key, Confirm: true})
		if err != nil || r.State != "retired" && r.State != "pending" {
			s.Items[i].State = "unknown"
			if e = save(); e != nil {
				return e
			}
			return acceptance.RequireAction(ctx, p.host, "cleanup-unconfirmed", "cleanup_outcome_unknown")
		}
		// pending is only accepted intent, never claimed as completed file deletion.
		s.Items[i].State = r.State
		if e = save(); e != nil {
			return e
		}
	}
	return nil
}
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) {
	return sdk.Form{Defaults: json.RawMessage(`{"active":false,"execute":false,"deleteFiles":false,"minRatio":"2","minSeedingHours":"168","batchSize":"10"}`), Content: []contracts.Component{
		acceptance.Field("VSwitch", "active", "启用定时预览"), acceptance.Field("VSwitch", "execute", "危险：实际删除任务（默认仅预览）"), acceptance.Field("VSwitch", "deleteFiles", "额外危险：请求删除数据文件"), acceptance.Field("VTextField", "minRatio", "最低分享率"), acceptance.Field("VTextField", "minSeedingHours", "最低做种小时"), acceptance.Field("VTextField", "batchSize", "每批上限（1–20）"), acceptance.Field("VTextField", "category", "分类（空表示不限）"), acceptance.Field("VTextField", "tag", "标签（空表示不限）"), acceptance.Field("VTextField", "siteId", "站点账户 ID（空表示不限）")}}, nil
}
func (p *Plugin) GetPage(ctx context.Context, _ sdk.PageRequest) (sdk.Page, error) {
	s, _, e := acceptance.Load[State](ctx, p.host, "cleanup-state")
	raw, _ := json.Marshal(s)
	return sdk.Page{Content: []contracts.Component{{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(string(raw))}}}}, e
}
