// Package checkin implements scheduled, host-authenticated sign-in. No credentials enter this process.
package checkin

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
	"github.com/dingding229/mediatool-plugin/sdk/host/notifications"
	"github.com/dingding229/mediatool-plugin/sdk/host/sites"
	"net/url"
	"strings"
	"sync"
	"time"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Config struct {
	Active        bool   `json:"active"`
	AccountID     string `json:"accountId"`
	Origin        string `json:"origin"`
	Path          string `json:"path"`
	PrecheckPath  string `json:"precheckPath,omitempty"`
	Method        string `json:"method"`
	Body          string `json:"body,omitempty"`
	SuccessMarker string `json:"successMarker"`
	Cron          string `json:"cron"`
}
type Result struct {
	Subject string `json:"subject"`
	Day     string `json:"day"`
	State   string `json:"state"`
}
type Plugin struct {
	mu      sync.Mutex
	host    *host.Client
	config  Config
	running bool
	now     func() time.Time
}

func New(s sdk.Services) sdk.Plugin { return &Plugin{host: s.Host, now: time.Now} }
func (p *Plugin) InitPlugin(_ context.Context, c sdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	cfg := Config{Path: "/attendance.php", Method: "GET", Cron: "0 8 * * *"}
	if contracts.Decode(c.Value, &cfg) != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	if cfg.Active {
		u, e := url.Parse(cfg.Origin)
		if e != nil || u.Scheme != "https" && u.Scheme != "http" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" || !contracts.OpaqueID(cfg.AccountID) || cfg.SuccessMarker == "" || len(cfg.SuccessMarker) > 200 || len(cfg.Body) > 16384 || cfg.Method != "GET" && cfg.Method != "POST" || cfg.Method == "GET" && cfg.Body != "" {
			return contracts.NewError(contracts.InvalidParams)
		}
		for _, path := range []string{cfg.Path, cfg.PrecheckPath} {
			if path != "" && (!strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") || strings.ContainsAny(path, "\\\r\n\x00") || len(path) > 2048) {
				return contracts.NewError(contracts.InvalidParams)
			}
		}
	}
	p.config = cfg
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
	return []sdk.Service{{ID: "checkin", Name: "受控站点签到（宿主时区）", Trigger: "cron", Cron: p.config.Cron, ScheduleSetting: "active", RequiredScopes: []string{"plugin.data.read", "plugin.data.write", "sites.request", "notifications.send"}, Handler: p.checkin}}, nil
}
func (p *Plugin) checkin(ctx context.Context, _ json.RawMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	day := p.now().UTC().Format("20060102")
	key := "checkin:" + p.config.AccountID + ":" + day
	old, rev, e := acceptance.Load[Result](ctx, p.host, "last-checkin")
	if e != nil {
		return e
	}
	subjectRaw, _ := json.Marshal(requestIdentity(p.config))
	subject := fmt.Sprintf("%x", sha256.Sum256(subjectRaw))
	if old.Day == day && old.Subject == subject {
		if old.State == "succeeded" {
			return nil
		}
		return contracts.NewError(contracts.NeedsAction)
	}
	request := b.SiteRequest{AccountID: p.config.AccountID, Origin: p.config.Origin, Path: p.config.Path, Method: p.config.Method, Body: p.config.Body, SuccessMarker: p.config.SuccessMarker}
	op := b.Operation{Key: key, Confirm: true}
	client := sites.New(p.host)
	result := b.SiteResponse{}
	if p.config.PrecheckPath != "" {
		probe := request
		probe.Path = p.config.PrecheckPath
		probe.Method = "GET"
		probe.Body = ""
		result, e = client.Request(ctx, probe, b.Operation{Key: key + ":probe", Confirm: true})
	}
	if e == nil && (p.config.PrecheckPath == "" || result.Status >= 200 && result.Status < 300) && !result.Matched && result.State != "needs_action" {
		result, e = client.Checkin(ctx, request, op)
	}
	state := "succeeded"
	if e != nil || !result.Matched || result.State != "matched" {
		state = "needs_action"
	}
	if _, saveErr := p.host.SaveData(ctx, "last-checkin", Result{Subject: subject, Day: day, State: state}, rev); saveErr != nil {
		return saveErr
	}
	if state == "needs_action" {
		code := "checkin_unconfirmed"
		if result.State == "needs_action" {
			code = "site_challenge"
		}
		return acceptance.RequireAction(ctx, p.host, key+":attention", code)
	}
	notice, e := notifications.New(p.host).Send(ctx, b.Notice{Title: "站点签到已确认", Text: "宿主已匹配管理员声明的成功标记；不包含站点正文或凭据。"}, b.Operation{Key: key + ":notice", Confirm: true})
	if e != nil || notice.State != b.NotificationPending {
		return acceptance.RequireAction(ctx, p.host, key+":notice-attention", "notification_unconfirmed")
	}
	return nil
}
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) {
	return sdk.Form{Content: []contracts.Component{
		acceptance.Field("VSwitch", "active", "定时签到", "启用定时签到（须先配置）"), acceptance.Field("VTextField", "accountId", "站点账号", "已配置站点账号 ID"), acceptance.Field("VTextField", "origin", "站点地址", "精确 Origin（非凭据）"), acceptance.Field("VTextField", "path", "签到页面", "签到相对路径"), acceptance.Field("VTextField", "precheckPath", "状态检查页面", "可选只读签到状态路径"), acceptance.Select("method", "请求方式", []string{"", "GET", "POST"}, map[string]string{"": "使用默认（GET）", "GET": "GET", "POST": "POST"}), acceptance.Field("VTextField", "body", "提交内容", "可选 POST 表单（禁止填凭据）"), acceptance.Field("VTextField", "successMarker", "成功标记", "成功标记（字面量，≤200 字节）"), acceptance.Field("VCronInput", "cron", "签到计划", "Cron（Asia/Shanghai，至少1分钟）")}, Defaults: json.RawMessage(`{"active":false,"path":"/attendance.php","method":"GET","cron":"0 8 * * *"}`)}, nil
}
func (p *Plugin) GetPage(ctx context.Context, _ sdk.PageRequest) (sdk.Page, error) {
	v, _, e := acceptance.Load[Result](ctx, p.host, "last-checkin")
	return sdk.Page{Content: []contracts.Component{{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(v.State + " " + v.Day)}}}}, e
}

// Schedule/active toggles do not change a request identity. Request edits cannot
// inherit another account's success; the host still fences same-day key conflicts.
func requestIdentity(c Config) b.SiteRequest {
	return b.SiteRequest{AccountID: c.AccountID, Origin: c.Origin, Path: c.Path, Method: c.Method, Body: c.Body, SuccessMarker: c.SuccessMarker}
}
