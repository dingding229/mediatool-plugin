package business_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	dl "github.com/dingding229/mediatool-plugin/sdk/host/downloads"
	sr "github.com/dingding229/mediatool-plugin/sdk/host/search"
	st "github.com/dingding229/mediatool-plugin/sdk/host/sites"
	ss "github.com/dingding229/mediatool-plugin/sdk/host/subscriptions"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

type caller struct{ h *plugintest.Host }

func (c caller) Call(ctx context.Context, m string, in, out any) error {
	if m != "host.call" {
		panic(m)
	}
	raw, e := json.Marshal(in)
	if e != nil {
		return e
	}
	v, rpc := c.h.Handle(ctx, raw)
	if rpc != nil {
		return rpc
	}
	raw, e = json.Marshal(v)
	if e != nil {
		return e
	}
	return json.Unmarshal(raw, out)
}
func TestEveryBusinessRPCScopeSuccessAndRevocation(t *testing.T) {
	ctx := context.Background()
	for _, m := range b.Methods() {
		t.Run(m.Name, func(t *testing.T) {
			mock := plugintest.NewHost(nil)
			client := host.New(caller{mock})
			var calls int
			if e := mock.OnBusiness(m.Name, func(context.Context, json.RawMessage) (any, *contracts.Error) {
				calls++
				return b.Receipt{ID: "fixture", State: "succeeded"}, nil
			}); e != nil {
				t.Fatal(e)
			}
			var rpc *contracts.Error
			op := b.Operation{Key: "fixture-operation", Confirm: true}
			if _, e := b.Call[b.Receipt](ctx, client, m.Name, struct{}{}, op); !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied || calls != 0 {
				t.Fatal("scope denial", e)
			}
			mock = plugintest.NewHost(append([]string{m.Scope}, b.RequiredScopes(m.Name)...))
			client = host.New(caller{mock})
			_ = mock.OnBusiness(m.Name, func(context.Context, json.RawMessage) (any, *contracts.Error) {
				calls++
				return b.Receipt{ID: "fixture", State: "succeeded"}, nil
			})
			out, e := b.Call[b.Receipt](ctx, client, m.Name, struct{}{}, op)
			if m.Supported {
				if e != nil || out.ID != "fixture" || calls != 1 {
					t.Fatal("success", e, out)
				}
			} else if !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported || calls != 0 {
				t.Fatal("unsupported", e)
			}
			if m.Supported {
				_ = mock.OnBusiness(m.Name, func(context.Context, json.RawMessage) (any, *contracts.Error) {
					return nil, contracts.NewError(contracts.Unsupported)
				})
				// A new logical operation observes changed backend capability; the old receipt remains immutable.
				op.Key = "fixture-unsupported-operation"
				if _, e = b.Call[b.Receipt](ctx, client, m.Name, struct{}{}, op); !errors.As(e, &rpc) || rpc.Code != contracts.Unsupported {
					t.Fatal("unsupported backend", e)
				}
			}
			mock.Revoke(m.Scope)
			if _, e = b.Call[b.Receipt](ctx, client, m.Name, struct{}{}, op); !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
				t.Fatal("revocation", e)
			}
		})
	}
}
func TestMethodGoldenAndBounds(t *testing.T) {
	plugintest.Golden(t, "testdata/methods.json", b.Methods())
	seen := map[string]bool{}
	for _, m := range b.Methods() {
		if seen[m.Name] || !contracts.KnownScope(m.Scope) {
			t.Fatal(m)
		}
		seen[m.Name] = true
	}
	for _, p := range []b.PageRequest{{Page: 0, Size: 10}, {Page: 1, Size: 201}, {Page: 1001, Size: 1}} {
		if p.Valid() {
			t.Fatal(p)
		}
	}
}
func TestTypedModuleClients(t *testing.T) {
	ctx := context.Background()
	m := plugintest.NewHost([]string{"downloads.read", "search.read", "sites.read", "subscriptions.read"})
	h := host.New(caller{m})
	for _, method := range []string{"downloads.list", "search.list", "sites.list", "subscriptions.list"} {
		_ = m.OnBusiness(method, func(context.Context, json.RawMessage) (any, *contracts.Error) {
			return map[string]any{"items": []any{}, "hasMore": false}, nil
		})
	}
	p := b.PageRequest{Page: 1, Size: 200}
	if _, e := dl.New(h).List(ctx, p); e != nil {
		t.Fatal(e)
	}
	if _, e := sr.New(h).List(ctx, p); e != nil {
		t.Fatal(e)
	}
	if _, e := st.New(h).List(ctx, p); e != nil {
		t.Fatal(e)
	}
	if _, e := ss.New(h).List(ctx, p); e != nil {
		t.Fatal(e)
	}
}
