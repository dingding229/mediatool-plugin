// Package refresh coalesces committed import events, using durable Host receipts for effects.
package refresh

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/dingding229/mediatool-plugin/plugins/internal/acceptance"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/mediaserver"
	"github.com/dingding229/mediatool-plugin/sdk/host/notifications"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Config struct {
	Active   bool    `json:"active"`
	Provider string  `json:"provider"`
	Interval Seconds `json:"intervalSeconds"`
}

const replayWindow = 7 * 24 * time.Hour

// Receipts retain 7d. Completed effects retain 8d, including 24h clock/scheduling margin.
const ledgerRetention = replayWindow + 24*time.Hour

type Notice struct {
	OccurredAt  time.Time `json:"occurredAt,omitempty"`
	CompletedAt int64     `json:"completedAt,omitempty"`
	ID          string    `json:"id"`
	AssetID     string    `json:"assetId"`
	Notified    bool      `json:"notified"`
}
type Pending struct {
	AssetID string `json:"assetId"`
	Receipt string `json:"receipt"`
}
type State struct {
	Events    []Notice  `json:"events"`
	Pending   []Pending `json:"pending"`
	LastFlush int64     `json:"lastFlush"`
}
type Plugin struct {
	mu      sync.Mutex
	host    *host.Client
	events  *events.Registry
	config  Config
	running bool
	now     func() time.Time
}

func New(s sdk.Services) sdk.Plugin { return &Plugin{host: s.Host, events: s.Events, now: time.Now} }
func (p *Plugin) InitPlugin(_ context.Context, c sdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v := Config{Provider: "builtin", Interval: 60}
	if contracts.Decode(c.Value, &v) != nil || v.Provider != "builtin" || v.Interval < 60 || v.Interval > 3600 {
		return contracts.NewError(contracts.InvalidParams)
	}
	p.config = v
	p.running = true
	return p.events.OnWindow(contracts.EventTransferComplete, replayWindow, p.onTransfer)
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
func (p *Plugin) onTransfer(ctx context.Context, e events.Envelope) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	var fact contracts.TransferCompleted
	if !contracts.OpaqueID(e.ID) || contracts.Decode(e.Payload, &fact) != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	// Legacy generic file operations have no media asset identity. Never guess one.
	if fact.AssetID == "" {
		return nil
	}
	if !contracts.OpaqueID(fact.AssetID) || fact.State != "succeeded" || fact.OperationID != "" && !contracts.OpaqueID(fact.OperationID) {
		return contracts.NewError(contracts.InvalidParams)
	}
	if e.OccurredAt.IsZero() || e.OccurredAt.After(p.now().Add(time.Minute)) {
		return contracts.NewError(contracts.InvalidParams)
	}
	if !e.OccurredAt.Add(replayWindow).After(p.now()) {
		return nil
	}
	s, rev, err := acceptance.Load[State](ctx, p.host, "refresh-state")
	if err != nil {
		return err
	}
	kept := s.Events[:0]
	for _, v := range s.Events {
		pending := false
		for _, q := range s.Pending {
			pending = pending || q.AssetID == v.AssetID
		}
		expired := !v.OccurredAt.IsZero() && !v.OccurredAt.Add(replayWindow).After(p.now())
		if !pending && v.Notified && v.CompletedAt > 0 && expired && p.now().Unix()-v.CompletedAt > int64(ledgerRetention/time.Second) {
			continue
		}
		kept = append(kept, v)
	}
	s.Events = kept
	index := -1
	for n, v := range s.Events {
		if v.ID == e.ID {
			index = n
			break
		}
	}
	if index < 0 {
		// Only completed, expired facts are evicted; pending/legacy/unknown never are.
		if len(s.Events) >= 200 {
			return acceptance.RequireAction(ctx, p.host, "refresh-ledger-full", "refresh_ledger_full")
		}
		index = len(s.Events)
		s.Events = append(s.Events, Notice{ID: e.ID, AssetID: fact.AssetID, OccurredAt: e.OccurredAt})
		found := false
		for _, v := range s.Pending {
			found = found || v.AssetID == fact.AssetID
		}
		if !found {
			s.Pending = append(s.Pending, Pending{fact.AssetID, e.ID})
		}
		next, e := p.host.SaveData(ctx, "refresh-state", s, rev)
		if e != nil {
			return e
		}
		rev = next.Revision
	}
	if s.Events[index].Notified {
		return nil
	}
	notice, err := notifications.New(p.host).Send(ctx, b.Notice{Title: "媒体入库已确认", Text: "资产 " + fact.AssetID + " 已加入插件待刷新队列；不代表媒体服务器刷新完成。"}, b.Operation{Key: "import-notice:" + e.ID, Confirm: true})
	if err != nil {
		return err
	}
	if notice.State != b.NotificationPending {
		return acceptance.RequireAction(ctx, p.host, "import-attention:"+e.ID, "notification_unconfirmed")
	}
	s.Events[index].Notified = true
	_, err = p.host.SaveData(ctx, "refresh-state", s, rev)
	return err
}
func (p *Plugin) GetService(context.Context) ([]sdk.Service, error) {
	return []sdk.Service{{ID: "flush", Name: "合并媒体刷新队列", Trigger: "interval", Interval: time.Minute, ScheduleSetting: "active", RequiredScopes: []string{"plugin.data.read", "plugin.data.write", "mediaserver.refresh", "notifications.send"}, Handler: p.flush}}, nil
}
func (p *Plugin) flush(ctx context.Context, _ json.RawMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	s, rev, e := acceptance.Load[State](ctx, p.host, "refresh-state")
	if e != nil {
		return e
	}
	if len(s.Pending) == 0 || p.now().Unix()-s.LastFlush < int64(p.config.Interval) {
		return nil
	}
	first := s.Pending[0]
	receipt, e := mediaserver.New(p.host).Refresh(ctx, b.ServerRefresh{Provider: p.config.Provider, AssetID: first.AssetID}, b.Operation{Key: "refresh:" + first.Receipt, Confirm: true})
	if e != nil || receipt.State != "requested" && receipt.State != "queued" {
		return acceptance.RequireAction(ctx, p.host, "refresh-attention:"+first.Receipt, "refresh_unconfirmed")
	}
	for i := range s.Events {
		if s.Events[i].AssetID == first.AssetID && s.Events[i].Notified {
			s.Events[i].CompletedAt = p.now().Unix()
		}
	}
	s.Pending = s.Pending[1:]
	s.LastFlush = p.now().Unix()
	_, e = p.host.SaveData(ctx, "refresh-state", s, rev)
	return e
}
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) {
	return sdk.Form{Content: []contracts.Component{acceptance.Field("VSwitch", "active", "启用入库通知及合并刷新"), acceptance.Select("provider", "媒体服务器", []string{"", "builtin"}, map[string]string{"": "使用默认（内置媒体库）", "builtin": "内置媒体库"}), intervalField()}, Defaults: json.RawMessage(`{"active":false,"provider":"builtin","intervalSeconds":"60"}`)}, nil
}
func (p *Plugin) GetPage(ctx context.Context, _ sdk.PageRequest) (sdk.Page, error) {
	s, _, e := acceptance.Load[State](ctx, p.host, "refresh-state")
	if e != nil {
		return sdk.Page{}, e
	}
	status, _, e := acceptance.Load[map[string]string](ctx, p.host, "status")
	return sdk.Page{Content: []contracts.Component{{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(fmt.Sprintf("待刷新 %d；事件 %d；%s %s", len(s.Pending), len(s.Events), status["state"], status["code"]))}}}}, e
}

func intervalField() contracts.Component {
	return acceptance.Number("intervalSeconds", "刷新间隔（秒）", 60, 3600, 1, true)
}

// The host whitelist renderer intentionally emits strings for VSelect/ElInput.
// Accept decimal form seconds and raw JSON integers; Init enforces 60..3600.
type Seconds int64

func (v *Seconds) UnmarshalJSON(raw []byte) error {
	text := strings.TrimSpace(string(raw))
	if len(text) > 0 && text[0] == '"' {
		if err := json.Unmarshal(raw, &text); err != nil {
			return err
		}
	}
	n, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return err
	}
	*v = Seconds(n)
	return nil
}
