// Package hello is a harmless lifecycle example: it increments its own KV counter.
package hello

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"github.com/dingding229/mediatool-plugin/sdk/host"
)

//go:embed manifest.json
var manifest string

func ManifestJSON() []byte { return []byte(manifest) }

type Plugin struct {
	host    *host.Client
	events  *events.Registry
	running bool
	config  map[string]json.RawMessage
	count   int
}

var _ sdk.Plugin = (*Plugin)(nil)
var _ sdk.FormProvider = (*Plugin)(nil)
var _ sdk.PageProvider = (*Plugin)(nil)

func New(s sdk.Services) sdk.Plugin { return &Plugin{host: s.Host, events: s.Events} }
func (p *Plugin) InitPlugin(ctx context.Context, c sdk.Config) error {
	if err := json.Unmarshal(c.Value, &p.config); err != nil {
		return err
	}
	revision, e := p.host.GetData(ctx, "starts", &p.count)
	var rpc *contracts.Error
	if e != nil && (!errors.As(e, &rpc) || rpc.Code != contracts.NotFound) {
		return e
	}
	if _, e = p.host.SaveData(ctx, "starts", p.count+1, revision); e != nil {
		return e
	}
	p.count++
	if e = p.events.On(events.EventType("library.organized"), p.onTransfer); e != nil {
		return e
	}
	if e = p.events.On(events.EventType("download.added"), p.onDownload); e != nil {
		return e
	}
	if e = p.events.OnChainPriority(events.ChainEventType("name.recognize"), 20, p.recognize); e != nil {
		return e
	}
	if e = p.acquisitionEvents(); e != nil {
		return e
	}
	if e = p.organizingEvents(); e != nil {
		return e
	}
	if e = p.discoveryEvents(); e != nil {
		return e
	}
	if e = p.platformEvents(); e != nil {
		return e
	}
	if e = p.workflowEvents(); e != nil {
		return e
	}
	p.running = true
	return nil
}
func (p *Plugin) GetState(context.Context) (bool, error)    { return p.running, nil }
func (p *Plugin) StopService(context.Context) error         { p.running = false; return nil }
func (p *Plugin) GetForm(context.Context) (sdk.Form, error) { return configurationForm(), nil }
func (p *Plugin) GetPage(context.Context, sdk.PageRequest) (sdk.Page, error) {
	b, _ := json.Marshal(p.count)
	return sdk.Page{Content: []contracts.Component{{Component: "Text", Key: "starts", Props: map[string]json.RawMessage{"count": b}}}}, nil
}
