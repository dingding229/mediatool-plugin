package hello

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"time"
)

func (p *Plugin) increment(ctx context.Context, key string) error {
	var n int
	revision, e := p.host.GetData(ctx, key, &n)
	var rpc *contracts.Error
	if e != nil && (!errors.As(e, &rpc) || rpc.Code != contracts.NotFound) {
		return e
	}
	_, e = p.host.SaveData(ctx, key, n+1, revision)
	return e
}
func (p *Plugin) GetAPI(context.Context) ([]sdk.API, error) {
	routes := []sdk.API{}
	for _, auth := range []string{"admin", "apikey", "bearer", "public"} {
		required := []string{"plugin.data.read"}
		methods := []string{"GET"}
		if auth == "admin" {
			methods = append(methods, "POST")
			required = append(required, "plugin.data.write")
		}
		routes = append(routes, sdk.API{Path: "/" + auth, Methods: methods, Auth: auth, Summary: "Hello counter", Description: "Only this instance's synthetic counter", RequiredScopes: required, Handler: func(ctx context.Context, q sdk.APIRequest) (sdk.APIResponse, error) {
			if q.Method == "POST" {
				if e := p.increment(ctx, "api_writes"); e != nil {
					return sdk.APIResponse{}, e
				}
			}
			var n int
			_, e := p.host.GetData(ctx, "ticks", &n)
			var rpc *contracts.Error
			if e != nil && (!errors.As(e, &rpc) || rpc.Code != contracts.NotFound) {
				return sdk.APIResponse{}, e
			}
			b, e := json.Marshal(map[string]any{"ticks": n, "instance": "hello-go"})
			return sdk.APIResponse{Status: 200, Headers: map[string]string{"Content-Type": "application/json"}, Body: b}, e
		}})
	}
	return routes, nil
}
func (p *Plugin) GetService(context.Context) ([]sdk.Service, error) {
	return []sdk.Service{{ID: "tick", Name: "计数示例", Trigger: "interval", Interval: time.Second, ScheduleSetting: "scheduled", Kwargs: json.RawMessage(`{}`), Handler: func(ctx context.Context, _ json.RawMessage) error { return p.increment(ctx, "ticks") }}}, nil
}
func (p *Plugin) GetCommand(context.Context) ([]sdk.Command, error) {
	return []sdk.Command{{Cmd: "hello", Desc: "增加实例计数", Category: "示例", Data: json.RawMessage(`{}`), Handler: func(ctx context.Context, _ json.RawMessage) (json.RawMessage, error) {
		return json.RawMessage(`{"ok":true}`), p.increment(ctx, "commands")
	}}}, nil
}
func (p *Plugin) GetActions(context.Context) ([]sdk.Action, error) {
	return []sdk.Action{{ID: "f4-check", Name: "F4 框架验收", RequiredScopes: []string{"plugin.data.read", "plugin.data.write"}, Kwargs: json.RawMessage(`{}`), Handler: p.f4Check}, {ID: "f3-read-check", Name: "F3 只读验收", RequiredScopes: []string{"plugin.data.read", "plugin.data.write", "downloads.read", "library.read", "mediaserver.read", "cache.read", "integrations.read", "jobs.read", "maintenance.read", "notifications.read", "settings.read", "system.read"}, Kwargs: json.RawMessage(`{}`), Handler: p.f3ReadCheck}, {ID: "hello", Name: "示例动作", RequiredScopes: []string{"plugin.data.read", "plugin.data.write"}, Kwargs: json.RawMessage(`{}`), Handler: func(ctx context.Context, q sdk.ActionContext) (sdk.ActionContext, bool, error) {
		q.Data["hello"] = json.RawMessage(`true`)
		return q, true, p.increment(ctx, "actions")
	}}}, nil
}
func (p *Plugin) onTransfer(ctx context.Context, e events.Envelope) error {
	// The durable delivery identity is also the key: duplicate delivery has no effect.
	var seen bool
	_, err := p.host.GetData(ctx, "event-"+e.ID, &seen)
	var rpc *contracts.Error
	if err == nil && seen {
		return nil
	}
	if err != nil && (!errors.As(err, &rpc) || rpc.Code != contracts.NotFound) {
		return err
	}
	_, err = p.host.SaveData(ctx, "event-"+e.ID, true, 0)
	return err
}
