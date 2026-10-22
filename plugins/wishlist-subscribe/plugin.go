// Package wishlist imports bounded public-list identities, never guesses a match.
package wishlist

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
	"github.com/dingding229/mediatool-plugin/sdk/host/metadata"
	"github.com/dingding229/mediatool-plugin/sdk/host/subscriptions"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Config struct {
	Active          bool               `json:"active"`
	Source          string             `json:"source"`
	ID              string             `json:"id"`
	AccountIDs      string             `json:"accountIds"`
	Create          bool               `json:"create"`
	IntervalMinutes acceptance.Numeric `json:"intervalMinutes"`
	PerRun          acceptance.Numeric `json:"perRun"`
}
type Pending struct {
	Entry      b.WishlistEntry `json:"entry"`
	Candidates []b.Media       `json:"candidates"`
}
type State struct {
	Page     int       `json:"page"`
	Offset   int       `json:"offset"`
	Pending  []Pending `json:"pending"`
	NextAt   int64     `json:"nextAt"`
	Failures int       `json:"failures"`
	Status   string    `json:"status"`
	Ensured  int       `json:"ensured"`
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
func accounts(s string) []string {
	out := []string{}
	for _, v := range strings.Split(s, ",") {
		if x := strings.TrimSpace(v); x != "" {
			out = append(out, x)
		}
	}
	return out
}
func (p *Plugin) InitPlugin(_ context.Context, c sdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v := Config{Source: "tmdb-list", IntervalMinutes: 60, PerRun: 10}
	if contracts.Decode(c.Value, &v) != nil || v.IntervalMinutes < 30 || v.IntervalMinutes > 1440 || math.Trunc(float64(v.IntervalMinutes)) != float64(v.IntervalMinutes) || v.PerRun < 1 || v.PerRun > 20 || math.Trunc(float64(v.PerRun)) != float64(v.PerRun) || v.Source != "tmdb-list" && v.Source != "douban-wish" && v.Source != "douban-subject" || len(v.ID) > 64 || v.Active && v.ID == "" {
		return contracts.NewError(contracts.InvalidParams)
	}
	ids := accounts(v.AccountIDs)
	if len(ids) > 10 || v.Active && len(ids) == 0 {
		return contracts.NewError(contracts.InvalidParams)
	}
	for _, id := range ids {
		if !contracts.OpaqueID(id) {
			return contracts.NewError(contracts.InvalidParams)
		}
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
	return []string{"plugin.data.read", "plugin.data.write", "metadata.read", "subscriptions.read", "subscriptions.write", "library.read"}
}
func (p *Plugin) GetService(context.Context) ([]sdk.Service, error) {
	return []sdk.Service{{ID: "sync", Name: "拉取公开想看并创建订阅", Trigger: "interval", Interval: 30 * time.Minute, ScheduleSetting: "active", RequiredScopes: scopes(), Handler: p.run}}, nil
}
func (p *Plugin) save(ctx context.Context, s State, rev int64) error {
	_, e := p.host.SaveData(ctx, "wishlist-state", s, rev)
	return e
}
func (p *Plugin) backoff(ctx context.Context, s State, rev int64, code string) error {
	s.Failures = min(s.Failures+1, 10)
	s.NextAt = p.now().Add(min(5*time.Minute*time.Duration(1<<min(s.Failures-1, 7)), 6*time.Hour)).Unix()
	s.Status = code
	return p.save(ctx, s, rev)
}
func sameTitle(a, b string) bool {
	return strings.EqualFold(strings.Join(strings.Fields(a), " "), strings.Join(strings.Fields(b), " "))
}
func (p *Plugin) matches(ctx context.Context, e b.WishlistEntry) ([]b.Media, error) {
	if e.TMDBID > 0 && (e.MediaType == "movie" || e.MediaType == "series") {
		return []b.Media{{ID: e.TMDBID, MediaType: e.MediaType, Title: e.Title, Year: e.Year}}, nil
	}
	out := []b.Media{}
	for _, kind := range []string{"movie", "series"} {
		if e.MediaType != "" && e.MediaType != kind {
			continue
		}
		r, err := metadata.New(p.host).Search(ctx, b.MetadataQuery{MediaType: kind, Query: e.Title})
		if err != nil {
			return nil, err
		}
		for _, m := range r {
			if sameTitle(m.Title, e.Title) && (e.Year == "" || m.Year == e.Year) {
				name := []rune(m.Title)
				if len(name) > 64 {
					name = name[:64]
				}
				out = append(out, b.Media{ID: m.ID, MediaType: m.MediaType, Title: string(name), Year: m.Year})
				if len(out) > 20 {
					return nil, contracts.NewError(contracts.PayloadTooLarge)
				}
			}
		}
	}
	return out, nil
}
func (p *Plugin) ensure(ctx context.Context, m b.Media) (string, error) {
	key := fmt.Sprintf("wishlist:%x", sha256.Sum256([]byte(m.MediaType+":"+strconv.Itoa(m.ID))))
	r, e := subscriptions.New(p.host).Ensure(ctx, b.SubscriptionInput{MediaType: m.MediaType, TMDBID: m.ID, Enabled: false, Source: "search", Rules: b.Rules{AccountIDs: accounts(p.config.AccountIDs)}}, b.Operation{Key: key, Confirm: true})
	if e != nil {
		return "", e
	}
	if r.State != "ensured" && r.State != "owned" {
		return "", contracts.NewError(contracts.ResultUnknown)
	}
	return r.State, nil
}
func (p *Plugin) run(ctx context.Context, _ json.RawMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	s, rev, e := acceptance.Load[State](ctx, p.host, "wishlist-state")
	if e != nil {
		return e
	}
	if p.now().Unix() < s.NextAt {
		return nil
	}
	if s.Page < 1 || s.Page > 1000 {
		s.Page = 1
	}
	if s.Offset < 0 || s.Offset > 200 {
		s.Offset = 0
	}
	page, e := metadata.New(p.host).Wishlist(ctx, b.WishlistRequest{Source: p.config.Source, ID: p.config.ID, Page: s.Page})
	if e != nil {
		return p.backoff(ctx, s, rev, "source_unavailable")
	}
	if page.State != "ok" {
		return p.backoff(ctx, s, rev, "source_"+page.State)
	}
	if len(page.Items) > 20 {
		return p.backoff(ctx, s, rev, "source_uncertain")
	}
	end := min(len(page.Items), s.Offset+int(p.config.PerRun))
	for s.Offset < end {
		entry := page.Items[s.Offset]
		if len(entry.Title) == 0 || len(entry.Title) > 200 || len(entry.SourceID) > 100 {
			return p.backoff(ctx, s, rev, "source_uncertain")
		}
		candidates, err := p.matches(ctx, entry)
		if err != nil {
			return p.backoff(ctx, s, rev, "lookup_unavailable")
		}
		if len(candidates) == 1 && p.config.Create && (entry.TMDBID > 0 || entry.Year != "") {
			if _, err = p.ensure(ctx, candidates[0]); err != nil {
				return p.backoff(ctx, s, rev, "ensure_unconfirmed")
			}
			s.Ensured++
		} else {
			found := false
			for _, v := range s.Pending {
				found = found || v.Entry.SourceID == entry.SourceID
			}
			if !found {
				if len(s.Pending) >= 200 {
					return p.backoff(ctx, s, rev, "pending_full")
				}
				next := append(append([]Pending{}, s.Pending...), Pending{Entry: entry, Candidates: candidates})
				raw, _ := json.Marshal(next)
				if len(raw) > 64<<10 {
					return p.backoff(ctx, s, rev, "pending_full")
				}
				s.Pending = next
			}
		}
		s.Offset++
	}
	if s.Offset >= len(page.Items) {
		s.Offset = 0
		if page.HasMore {
			s.Page++
		} else {
			s.Page = 1
		}
	}
	s.Failures = 0
	s.Status = "ok"
	s.NextAt = p.now().Add(time.Duration(p.config.IntervalMinutes) * time.Minute).Unix()
	return p.save(ctx, s, rev)
}
func (p *Plugin) GetActions(context.Context) ([]sdk.Action, error) {
	return []sdk.Action{{ID: "resolve", Name: "危险：确认候选并创建停用订阅", Kwargs: json.RawMessage(`{"sourceId":"","tmdbId":0,"mediaType":"movie"}`), RequiredScopes: scopes(), Handler: p.resolve}}, nil
}
func (p *Plugin) resolve(ctx context.Context, q sdk.ActionContext) (sdk.ActionContext, bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	var id, kind string
	var number int
	if json.Unmarshal(q.Data["sourceId"], &id) != nil || json.Unmarshal(q.Data["mediaType"], &kind) != nil || json.Unmarshal(q.Data["tmdbId"], &number) != nil {
		return q, false, contracts.NewError(contracts.InvalidParams)
	}
	s, rev, e := acceptance.Load[State](ctx, p.host, "wishlist-state")
	if e != nil {
		return q, false, e
	}
	for i, v := range s.Pending {
		if v.Entry.SourceID != id {
			continue
		}
		for _, m := range v.Candidates {
			if m.ID == number && m.MediaType == kind {
				state, e := p.ensure(ctx, m)
				if e != nil {
					return q, false, e
				}
				s.Pending = append(s.Pending[:i], s.Pending[i+1:]...)
				e = p.save(ctx, s, rev)
				return sdk.ActionContext{Data: map[string]json.RawMessage{"state": acceptance.Prop(state)}}, true, e
			}
		}
	}
	return q, false, contracts.NewError(contracts.InvalidParams)
}
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) {
	source := acceptance.Select("source", "公开来源", []string{"tmdb-list", "douban-wish", "douban-subject"}, map[string]string{"tmdb-list": "TMDB 列表", "douban-wish": "豆瓣公开想看", "douban-subject": "豆瓣公开条目"}, "不读取登录态收藏")
	return sdk.Form{Defaults: json.RawMessage(`{"active":false,"create":false,"source":"tmdb-list","intervalMinutes":"60","perRun":"10"}`), Content: []contracts.Component{acceptance.Field("VSwitch", "active", "启用定时拉取"), acceptance.Field("VSwitch", "create", "创建订阅", "危险：唯一匹配时创建停用订阅（默认只列候选）"), source, acceptance.Field("VTextField", "id", "来源标识", "列表 ID、豆瓣公开用户 ID 或条目 ID"), acceptance.Field("VTextField", "accountIds", "站点账号", "已配置 PT 账户 ID，逗号分隔（不立即搜索/下载）"), acceptance.Number("intervalMinutes", "同步间隔（分钟）", 30, 1440, 1, true, "同步间隔分钟（30–1440）"), acceptance.Number("perRun", "每轮条目上限", 1, 20, 1, true, "每轮条目（1–20）")}}, nil
}
func (p *Plugin) GetPage(ctx context.Context, _ sdk.PageRequest) (sdk.Page, error) {
	s, _, e := acceptance.Load[State](ctx, p.host, "wishlist-state")
	nodes := []contracts.Component{{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(fmt.Sprintf("状态 %s；待确认 %d；确保次数 %d；仅显示前5条，处理后刷新", s.Status, len(s.Pending), s.Ensured))}}}
	for i, v := range s.Pending {
		if i >= 5 {
			break
		}
		nodes = append(nodes, contracts.Component{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(v.Entry.SourceID + " / " + v.Entry.Title + " / " + v.Entry.Year)}})
		for _, m := range v.Candidates {
			nodes = append(nodes, contracts.Component{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(fmt.Sprintf("候选 tmdbId=%d mediaType=%s title=%s year=%s", m.ID, m.MediaType, m.Title, m.Year))}})
		}
	}
	return sdk.Page{Content: nodes}, e
}
