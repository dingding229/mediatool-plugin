// Package watch polls an explicit root-relative directory and submits core file plans.
// No worker filesystem access, recursive unbounded scans, moves or renames occur here.
package watch

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
	"github.com/dingding229/mediatool-plugin/sdk/host/files"
	"math"
	"path"
	"strings"
	"sync"
	"time"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Config struct {
	Active       bool               `json:"active"`
	RootID       string             `json:"rootId"`
	Directory    string             `json:"directory"`
	Destination  string             `json:"destination"`
	QuietSeconds acceptance.Numeric `json:"quietSeconds"`
}
type Observation struct {
	File        b.MediaFile `json:"file"`
	Key         string      `json:"key"`
	StableSince int64       `json:"stableSince"`
	Plan        *b.FilePlan `json:"plan,omitempty"`
	Applying    bool        `json:"applying"`
}
type State struct {
	Items     []Observation `json:"items"`
	Submitted int64         `json:"submitted"`
}
type Plugin struct {
	mu      sync.Mutex
	host    *host.Client
	config  Config
	running bool
	now     func() time.Time
}

func New(s sdk.Services) sdk.Plugin { return &Plugin{host: s.Host, now: time.Now} }
func relative(s string) bool {
	return s != "" && len(s) <= 1024 && path.Clean(s) == s && !strings.HasPrefix(s, "/") && s != ".." && !strings.HasPrefix(s, "../") && !strings.ContainsAny(s, "\\\x00\r\n")
}
func (p *Plugin) InitPlugin(_ context.Context, c sdk.Config) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v := Config{QuietSeconds: 60}
	if contracts.Decode(c.Value, &v) != nil || v.QuietSeconds < 30 || v.QuietSeconds > 3600 || math.Trunc(float64(v.QuietSeconds)) != float64(v.QuietSeconds) || v.Active && (!contracts.OpaqueID(v.RootID) || !relative(v.Directory) || !relative(v.Destination) || v.Directory == v.Destination || strings.HasPrefix(v.Destination, v.Directory+"/")) {
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
	return []sdk.Service{{ID: "scan", Name: "稳定文件补扫与宿主整理", Trigger: "interval", Interval: 30 * time.Second, ScheduleSetting: "active", RequiredScopes: []string{"plugin.data.read", "plugin.data.write", "fs.media:read", "fs.media:write"}, Handler: p.run}}, nil
}
func media(name string) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".mkv", ".mp4", ".avi", ".m4v", ".ts", ".m2ts":
		return true
	}
	return false
}
func (p *Plugin) run(ctx context.Context, _ json.RawMessage) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running || !p.config.Active {
		return nil
	}
	s, rev, e := acceptance.Load[State](ctx, p.host, "watch-state")
	if e != nil {
		return e
	}
	client := files.New(p.host)
	// One bounded directory (not a recursive tree). Overflow is visible, never silently dropped.
	page, e := client.List(ctx, b.FileQuery{RootID: p.config.RootID, Path: p.config.Directory, PageRequest: b.PageRequest{Page: 1, Size: 200}})
	if e != nil {
		return e
	}
	if page.HasMore || len(page.Items) > 200 {
		return acceptance.RequireAction(ctx, p.host, "watch-capacity", "watch_queue_full")
	}
	for i := range s.Items {
		v := &s.Items[i]
		if !v.Applying {
			continue
		}
		if v.Plan == nil {
			return acceptance.RequireAction(ctx, p.host, "watch-unknown", "watch_outcome_unknown")
		}
		r, err := client.Get(ctx, b.ID{ID: v.Plan.ID})
		if err != nil || r.State != "succeeded" {
			return acceptance.RequireAction(ctx, p.host, "watch-unknown", "watch_outcome_unknown")
		}
		v.Plan = &r
		v.Applying = false
		s.Submitted++
	}
	old := map[string]Observation{}
	for _, v := range s.Items {
		old[v.File.Path] = v
	}
	next := []Observation{}
	for _, f := range page.Items {
		if f.Unsupported || f.Directory || f.ID == "" || !media(f.Name) {
			continue
		}
		if f.Modified <= 0 {
			return contracts.NewError(contracts.Unsupported)
		}
		if !relative(f.Path) || path.Dir(f.Path) != p.config.Directory {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		key := fmt.Sprintf("watch:%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%s:%d:%d", p.config.RootID, p.config.Destination, f.Path, f.ID, f.Size, f.Modified))))
		v, ok := old[f.Path]
		if !ok || v.Key != key {
			v = Observation{File: f, Key: key, StableSince: p.now().Unix()}
		}
		next = append(next, v)
	}
	s.Items = next
	save := func() error {
		v, err := p.host.SaveData(ctx, "watch-state", s, rev)
		if err == nil {
			rev = v.Revision
		}
		return err
	}
	if e = save(); e != nil {
		return e
	}
	processed := 0
	for i := range s.Items {
		v := &s.Items[i]
		if processed >= 10 || p.now().Unix()-v.StableSince < int64(p.config.QuietSeconds) || v.Plan != nil && v.Plan.State == "succeeded" {
			continue
		}
		processed++
		if v.Plan == nil {
			size, mtime := v.File.Size, v.File.Modified
			r, err := client.Preview(ctx, b.FileInput{RootID: p.config.RootID, Operation: "organize", Source: v.File.Path, Destination: p.config.Destination, ExpectedFileID: v.File.ID, ExpectedSize: &size, ExpectedModified: &mtime}, b.Operation{Key: v.Key + ":preview", Confirm: true})
			if err != nil {
				return err
			}
			if r.State != "preview" || len(r.Digest) != 64 {
				return acceptance.RequireAction(ctx, p.host, "watch-preview", "watch_preview_unconfirmed")
			}
			v.Plan = &r
			if e = save(); e != nil {
				return e
			}
		}
		v.Applying = true
		if e = save(); e != nil {
			return e
		}
		r, err := client.Apply(ctx, b.Approval{ID: v.Plan.ID, Digest: v.Plan.Digest}, b.Operation{Key: v.Key + ":apply", Confirm: true})
		if err != nil || r.State != "succeeded" {
			return acceptance.RequireAction(ctx, p.host, "watch-apply", "watch_outcome_unknown")
		}
		v.Plan = &r
		v.Applying = false
		s.Submitted++
		if e = save(); e != nil {
			return e
		}
	}
	return nil
}
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) {
	return sdk.Form{Defaults: json.RawMessage(`{"active":false,"directory":"incoming","destination":"organized","quietSeconds":"60"}`), Content: []contracts.Component{acceptance.Field("VSwitch", "active", "危险：启用宿主预览/确认整理"), acceptance.Field("VTextField", "rootId", "已授权可写媒体根 ID"), acceptance.Field("VTextField", "directory", "监控相对目录（仅本层，最多 200 项）"), acceptance.Field("VTextField", "destination", "同一根内目标相对目录"), acceptance.Field("VTextField", "quietSeconds", "大小与 mtime 静默期（30–3600 秒）")}}, nil
}
func (p *Plugin) GetPage(ctx context.Context, _ sdk.PageRequest) (sdk.Page, error) {
	s, _, e := acceptance.Load[State](ctx, p.host, "watch-state")
	return sdk.Page{Content: []contracts.Component{{Component: "Text", Props: map[string]json.RawMessage{"text": acceptance.Prop(fmt.Sprintf("待观测 %d；已由宿主确认 %d", len(s.Items), s.Submitted))}}}}, e
}
