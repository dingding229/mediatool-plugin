package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"slices"
	"sort"
	"time"
)

type API struct {
	Path                       string
	Methods                    []string
	Auth, Summary, Description string
	RequiredScopes             []string
	Handler                    func(context.Context, APIRequest) (APIResponse, error)
}
type Service struct {
	ID, Name, Trigger, Cron string
	Interval                time.Duration
	At                      time.Time
	ScheduleSetting         string
	Kwargs                  json.RawMessage
	RequiredScopes          []string
	Handler                 func(context.Context, json.RawMessage) error
}
type Command struct {
	Cmd, Desc, Category string
	Data                json.RawMessage
	RequiredScopes      []string
	Handler             func(context.Context, json.RawMessage) (json.RawMessage, error)
}
type Action struct {
	ID, Name       string
	Kwargs         json.RawMessage
	RequiredScopes []string
	Handler        func(context.Context, ActionContext) (ActionContext, bool, error)
}
type Module struct {
	ID             string
	Priority       int
	RequiredScopes []string
	Handler        func(context.Context, json.RawMessage) (json.RawMessage, error)
}
type AuthProvider struct {
	ID, Name, Issuer, Audience, Origin string
	Handler                            func(context.Context, contracts.AuthRequest) (contracts.AuthAssertion, error)
}
type AuthProviderRegistry interface {
	GetAuthProviders(context.Context) ([]AuthProvider, error)
}
type bound struct {
	descriptor contracts.Descriptor
	call       func(context.Context, json.RawMessage) (any, error)
}

func (r *Runtime) bind(ctx context.Context) error {
	r.bindings = map[string]bound{}
	add := func(d contracts.Descriptor, h func(context.Context, json.RawMessage) (any, error)) error {
		d.Handler = contracts.HandlerID(d.Kind, d.ID)
		if d.RequiredScopes == nil {
			d.RequiredScopes = []string{}
		}
		if h == nil || d.ID == "" || len(d.ID) > 200 || r.bindings[d.Handler].call != nil || len(r.bindings) >= 128 || d.Priority < -1000 || d.Priority > 1000 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, scope := range d.RequiredScopes {
			if !slices.Contains(r.init.Capabilities, scope) {
				return contracts.NewError(contracts.PermissionDenied)
			}
		}
		r.bindings[d.Handler] = bound{d, h}
		return nil
	}
	if p, ok := r.plugin.(APIProvider); ok {
		list, e := p.GetAPI(ctx)
		if e != nil || len(list) > 64 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, a := range list {
			var route *contracts.APIRoute
			for i := range r.manifest.Extensions.API {
				v := &r.manifest.Extensions.API[i]
				if v.Path == a.Path {
					route = v
					break
				}
			}
			if route == nil || a.Handler == nil || a.Auth != route.Auth || !slices.Equal(a.Methods, route.Methods) || !slices.Equal(a.RequiredScopes, route.ExecutionScopes) || len(a.Description) > 2000 {
				return contracts.NewError(contracts.InvalidExtensionResult)
			}
			d := contracts.Descriptor{Kind: "api", ID: route.ID, Path: a.Path, Methods: a.Methods, Auth: a.Auth, Summary: a.Summary, Description: a.Description, RequiredScopes: append([]string{"api.register"}, a.RequiredScopes...)}
			if a.Auth == "public" {
				d.RequiredScopes = append(d.RequiredScopes, "api.public")
				for _, m := range a.Methods {
					if m != "GET" {
						d.RequiredScopes = append(d.RequiredScopes, "api.public.write")
						break
					}
				}
			}
			if e = add(d, func(ctx context.Context, b json.RawMessage) (any, error) {
				var q APIRequest
				if contracts.Decode(b, &q) != nil {
					return nil, contracts.NewError(contracts.InvalidParams)
				}
				return a.Handler(ctx, q)
			}); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(ServiceProvider); ok {
		list, e := p.GetService(ctx)
		if e != nil || len(list) > 64 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, v := range list {
			if v.Handler == nil || !slices.Contains(r.manifest.Extensions.Services, v.ID) {
				return contracts.NewError(contracts.InvalidExtensionResult)
			}
			d := contracts.Descriptor{Kind: "service", ID: v.ID, Name: v.Name, Trigger: v.Trigger, Cron: v.Cron, IntervalMS: v.Interval.Milliseconds(), ScheduleSetting: v.ScheduleSetting, Kwargs: v.Kwargs, RequiredScopes: append([]string{"services.register"}, v.RequiredScopes...)}
			if !v.At.IsZero() {
				d.At = v.At.UTC().Format(time.RFC3339Nano)
			}
			if e = add(d, func(ctx context.Context, b json.RawMessage) (any, error) { return struct{}{}, v.Handler(ctx, b) }); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(CommandProvider); ok {
		list, e := p.GetCommand(ctx)
		if e != nil || len(list) > 64 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, v := range list {
			if v.Handler == nil || !slices.Contains(r.manifest.Extensions.Commands, v.Cmd) {
				return contracts.NewError(contracts.InvalidExtensionResult)
			}
			if e = add(contracts.Descriptor{Kind: "command", ID: v.Cmd, Cmd: v.Cmd, Description: v.Desc, Category: v.Category, Data: v.Data, RequiredScopes: append([]string{"commands.register"}, v.RequiredScopes...)}, func(ctx context.Context, b json.RawMessage) (any, error) { return v.Handler(ctx, b) }); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(ActionProvider); ok {
		list, e := p.GetActions(ctx)
		if e != nil || len(list) > 64 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, v := range list {
			if v.Handler == nil || !slices.Contains(r.manifest.Extensions.Actions, v.ID) {
				return contracts.NewError(contracts.InvalidExtensionResult)
			}
			if e = add(contracts.Descriptor{Kind: "action", ID: v.ID, Name: v.Name, Kwargs: v.Kwargs, RequiredScopes: append([]string{"actions.register"}, v.RequiredScopes...)}, func(ctx context.Context, b json.RawMessage) (any, error) {
				var q ActionContext
				if contracts.Decode(b, &q) != nil || q.Data == nil {
					return nil, contracts.NewError(contracts.InvalidParams)
				}
				out, next, e := v.Handler(ctx, q)
				return contracts.ActionResult{Context: out, Continue: next}, e
			}); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(AuthProviderRegistry); ok {
		list, e := p.GetAuthProviders(ctx)
		if e != nil || len(list) != 1 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, v := range list {
			if v.Handler == nil || !slices.Contains(r.manifest.Extensions.AuthProviders, v.ID) {
				return contracts.NewError(contracts.InvalidExtensionResult)
			}
			d := contracts.Descriptor{Kind: "chain", ID: "auth.verification", Name: v.Name, ProviderID: v.ID, Issuer: v.Issuer, Audience: v.Audience, Origin: v.Origin, RequiredScopes: []string{"auth.provider", "events.chain.auth.verification"}}
			if e = add(d, func(ctx context.Context, b json.RawMessage) (any, error) {
				var envelope events.ChainEnvelope
				if contracts.Decode(b, &envelope) != nil || envelope.Type != contracts.ChainAuthVerification {
					return nil, contracts.NewError(contracts.InvalidParams)
				}
				var q contracts.AuthRequest
				if contracts.Decode(envelope.Payload, &q) != nil || q.ProviderID != v.ID || len(q.Challenge) > 200 {
					return nil, contracts.NewError(contracts.InvalidParams)
				}
				return v.Handler(ctx, q)
			}); e != nil {
				return e
			}
		}
	}
	if p, ok := r.plugin.(ModuleProvider); ok {
		list, e := p.GetModule(ctx)
		if e != nil || len(list) > 64 {
			return contracts.NewError(contracts.InvalidExtensionResult)
		}
		for _, v := range list {
			if v.Handler == nil || !slices.Contains(r.manifest.Extensions.Modules, v.ID) {
				return contracts.NewError(contracts.InvalidExtensionResult)
			}
			if e = add(contracts.Descriptor{Kind: "module", ID: v.ID, Priority: v.Priority, RequiredScopes: append([]string{"module." + v.ID}, v.RequiredScopes...)}, func(ctx context.Context, b json.RawMessage) (any, error) { return v.Handler(ctx, b) }); e != nil {
				return e
			}
		}
	}
	for _, v := range r.events.Descriptors() {
		if v.Kind == "chain" {
			typ := events.ChainEventType(v.Type)
			if e := add(contracts.Descriptor{Kind: "chain", ID: v.Type, Priority: v.Priority, RequiredScopes: []string{"events.chain." + v.Type}}, func(ctx context.Context, b json.RawMessage) (any, error) {
				var q events.ChainEnvelope
				if contracts.Decode(b, &q) != nil || q.Type != typ {
					return nil, contracts.NewError(contracts.InvalidParams)
				}
				return r.events.DispatchChain(ctx, q)
			}); e != nil {
				return e
			}
			continue
		}
		typ := events.EventType(v.Type)
		if e := add(contracts.Descriptor{Kind: "event", ID: v.Type, RequiredScopes: []string{"events.subscribe." + v.Type}}, func(ctx context.Context, b json.RawMessage) (any, error) {
			var q events.Envelope
			if contracts.Decode(b, &q) != nil || q.Type != typ {
				return nil, contracts.NewError(contracts.InvalidParams)
			}
			return struct{}{}, r.events.Dispatch(ctx, q)
		}); e != nil {
			return e
		}
	}
	if e := r.bindUI(ctx, add); e != nil {
		return e
	}
	return contracts.ValidateSnapshot(r.manifest, r.init.Capabilities, r.descriptors(""))
}
func (r *Runtime) descriptors(kind string) contracts.Snapshot {
	out := contracts.Snapshot{Version: contracts.DescriptorVersion, Descriptors: []contracts.Descriptor{}}
	for _, b := range r.bindings {
		if kind == "" || b.descriptor.Kind == kind {
			out.Descriptors = append(out.Descriptors, b.descriptor)
		}
	}
	sort.Slice(out.Descriptors, func(i, j int) bool {
		a, b := out.Descriptors[i], out.Descriptors[j]
		if a.Priority != b.Priority {
			return a.Priority > b.Priority
		}
		return a.Handler < b.Handler
	})
	return out
}
func (r *Runtime) invoke(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
	var q contracts.Invocation
	if contracts.Decode(b, &q) != nil || contracts.ValidateJSON(q.Args) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	h, ok := r.bindings[q.Handler]
	if !ok {
		return nil, contracts.NewError(contracts.NotFound)
	}
	out, e := h.call(ctx, q.Args)
	if e != nil {
		var code *contracts.Error
		if errors.As(e, &code) {
			return nil, contracts.SafeError(code)
		}
		return nil, contracts.NewError(contracts.NeedsAction)
	}
	raw, e := json.Marshal(out)
	if e != nil || len(raw) > 128<<10 || contracts.ValidateJSON(raw) != nil {
		return nil, contracts.NewError(contracts.InvalidExtensionResult)
	}
	return out, nil
}
