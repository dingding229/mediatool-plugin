package sdk_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/registry"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"github.com/dingding229/mediatool-plugin/sdk/schema"
	"sync/atomic"
	"testing"
	"time"
)

func fixture(t *testing.T) (registry.Definition, *plugintest.Host, sdk.Identity, contracts.Initialize) {
	t.Helper()
	d, e := registry.Get("hello-go")
	if e != nil {
		t.Fatal(e)
	}
	scopes := []string{}
	for _, p := range d.Manifest.Permissions {
		scopes = append(scopes, p.Scope)
	}
	id := sdk.Identity{InstanceID: "test-instance", DefinitionFingerprint: "test-definition", CatalogBuild: "test-build"}
	init := contracts.Initialize{ContractVersion: 2, InstanceID: id.InstanceID, Generation: "synthetic-generation", ConfigRevision: 1, DefinitionFingerprint: id.DefinitionFingerprint, CatalogBuild: id.CatalogBuild, Capabilities: scopes}
	return d, plugintest.NewHost(scopes), id, init
}
func start(t *testing.T, d registry.Definition, m *plugintest.Host, id sdk.Identity) *plugintest.Harness {
	t.Helper()
	h, e := plugintest.Start(context.Background(), d.Manifest, d.Factory, m, id)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() {
		if e := h.Close(); e != nil {
			t.Error(e)
		}
	})
	return h
}
func rpcCode(t *testing.T, e error, c contracts.Code) {
	t.Helper()
	var x *contracts.Error
	if !errors.As(e, &x) || x.Code != c {
		t.Fatalf("got %v expected %d", e, c)
	}
}
func TestHelloLifecycleGolden(t *testing.T) {
	d, m, id, q := fixture(t)
	h := start(t, d, m, id)
	ctx := context.Background()
	rpcCode(t, h.Call(ctx, "plugin.get_state", nil, nil), contracts.Unavailable)
	init, e := h.Initialize(ctx, q)
	if e != nil {
		t.Fatal(e)
	}
	if init.SchemaFingerprint != schema.Fingerprint {
		t.Fatal("schema mismatch")
	}
	if _, e = h.Initialize(ctx, q); e != nil {
		t.Fatal("same initialization should be idempotent", e)
	}
	if e = h.Init(ctx, contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}); e != nil {
		t.Fatal(e)
	}
	var state bool
	var form sdk.Form
	var page sdk.Page
	for _, call := range []struct {
		method    string
		args, out any
	}{{"plugin.get_state", nil, &state}, {"plugin.get_form", nil, &form}, {"plugin.get_page", sdk.PageRequest{Key: "overview"}, &page}} {
		if e = h.Call(ctx, call.method, call.args, call.out); e != nil {
			t.Fatal(e)
		}
	}
	plugintest.Golden(t, "testdata/hello-lifecycle.json", struct {
		Initialize contracts.Initialized `json:"initialize"`
		State      bool                  `json:"state"`
		Form       sdk.Form              `json:"form"`
		Page       sdk.Page              `json:"page"`
		Calls      []string              `json:"hostCalls"`
	}{init, state, form, page, m.Calls()})
	rpcCode(t, h.Init(ctx, contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}), contracts.Conflict)
	rpcCode(t, h.Call(ctx, "plugin.get_page", sdk.PageRequest{Key: "private"}, nil), contracts.NotFound)
	if e = h.Call(ctx, "plugin.get_api", nil, nil); e != nil {
		t.Fatal(e)
	}
	rpcCode(t, h.Call(ctx, "plugin.invoke", nil, nil), contracts.InvalidParams)
	for range 2 {
		if e = h.Call(ctx, "plugin.stop_service", nil, nil); e != nil {
			t.Fatal(e)
		}
	}
	rpcCode(t, h.Call(ctx, "plugin.get_state", nil, nil), contracts.NeedsAction)
}
func TestRejectedHandshakeNeverConstructs(t *testing.T) {
	for _, name := range []string{"version", "build", "definition", "instance", "missing-scope", "extra-scope", "duplicate-scope", "generation"} {
		t.Run(name, func(t *testing.T) {
			d, m, id, q := fixture(t)
			var created atomic.Int32
			original := d.Factory
			d.Factory = func(s sdk.Services) sdk.Plugin { created.Add(1); return original(s) }
			h := start(t, d, m, id)
			switch name {
			case "version":
				q.ContractVersion = 1
			case "build":
				q.CatalogBuild = "wrong"
			case "definition":
				q.DefinitionFingerprint = "wrong"
			case "instance":
				q.InstanceID = "other"
			case "missing-scope":
				q.Capabilities = q.Capabilities[1:]
			case "extra-scope":
				q.Capabilities = append(q.Capabilities, "system.exec")
			case "duplicate-scope":
				q.Capabilities = append(q.Capabilities, q.Capabilities[0])
			case "generation":
				q.Generation = ""
			}
			if _, e := h.Initialize(context.Background(), q); e == nil {
				t.Fatal("rejected handshake accepted")
			}
			if created.Load() != 0 {
				t.Fatal("constructor ran before authorization")
			}
		})
	}
}
func TestFactoryDeferredAndDisconnectStops(t *testing.T) {
	d, m, id, q := fixture(t)
	var created, stopped atomic.Int32
	d.Manifest.Extensions = contracts.Extensions{}
	d.Factory = func(s sdk.Services) sdk.Plugin { created.Add(1); return &minimal{stopped: &stopped} }
	h := start(t, d, m, id)
	if created.Load() != 0 {
		t.Fatal("catalog constructed plugin")
	}
	if _, e := h.Initialize(context.Background(), q); e != nil {
		t.Fatal(e)
	}
	if created.Load() != 0 {
		t.Fatal("initialize constructed plugin")
	}
	if e := h.Init(context.Background(), contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}); e != nil {
		t.Fatal(e)
	}
	if e := h.Close(); e != nil {
		t.Fatal(e)
	}
	if stopped.Load() != 1 {
		t.Fatal("missing disconnect cleanup")
	}
}

type minimal struct{ stopped *atomic.Int32 }

func (*minimal) InitPlugin(context.Context, sdk.Config) error { return nil }
func (*minimal) GetState(context.Context) (bool, error)       { return true, nil }
func (m *minimal) StopService(context.Context) error          { m.stopped.Add(1); return nil }
func TestPanicIsSanitizedAndFailedInitCannotRun(t *testing.T) {
	d, m, id, q := fixture(t)
	d.Factory = func(s sdk.Services) sdk.Plugin { panic("synthetic secret?token=hidden") }
	h := start(t, d, m, id)
	if _, e := h.Initialize(context.Background(), q); e != nil {
		t.Fatal(e)
	}
	rpcCode(t, h.Init(context.Background(), contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}), contracts.InternalError)
	rpcCode(t, h.Call(context.Background(), "plugin.get_state", nil, nil), contracts.NeedsAction)
}
func TestBadConfigAndLateHandshake(t *testing.T) {
	d, m, id, q := fixture(t)
	h := start(t, d, m, id)
	ctx := context.Background()
	if _, e := h.Initialize(ctx, q); e != nil {
		t.Fatal(e)
	}
	rpcCode(t, h.Init(ctx, contracts.Config{Value: json.RawMessage(`{}`), Revision: 2}), contracts.InvalidParams)
	q.Generation = "another"
	_, e := h.Initialize(ctx, q)
	rpcCode(t, e, contracts.Conflict)
}
func TestHarnessCancellationIsBounded(t *testing.T) {
	d, m, id, _ := fixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	h, e := plugintest.Start(ctx, d.Manifest, d.Factory, m, id)
	if e != nil {
		t.Fatal(e)
	}
	cancel()
	done := make(chan error, 1)
	go func() { done <- h.Close() }()
	select {
	case e = <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("harness cancellation blocked")
	}
}
