package host_test

import (
	"context"
	"errors"
	"github.com/dingding229/mediatool-plugin/registry"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"strings"
	"testing"
	"time"
)

func TestTypedHelpersRevisionScopesAndIsolation(t *testing.T) {
	d, e := registry.Get("hello-go")
	if e != nil {
		t.Fatal(e)
	}
	scopes := []string{"plugin.config.read", "plugin.config.write", "plugin.data.read", "plugin.data.write", "fs.data", "notifications.send"}
	m := plugintest.NewHost(scopes)
	h, e := plugintest.Start(context.Background(), d.Manifest, d.Factory, m, sdk.Identity{InstanceID: "test", DefinitionFingerprint: "def", CatalogBuild: "build"})
	if e != nil {
		t.Fatal(e)
	}
	defer h.Close()
	c := h.Client
	ctx := context.Background()
	var config map[string]string
	if rev, e := c.GetConfig(ctx, &config); e != nil || rev != 1 {
		t.Fatal(e, rev)
	}
	if r, e := c.UpdateConfig(ctx, map[string]string{"greeting": "hi"}, 1); e != nil || r.Revision != 2 {
		t.Fatal(e, r)
	}
	if _, e = c.UpdateConfig(ctx, map[string]string{}, 1); e == nil {
		t.Fatal("stale config revision accepted")
	}
	if r, e := c.SaveData(ctx, "counter", 4, 0); e != nil || r.Revision != 1 {
		t.Fatal(e, r)
	}
	var out int
	if rev, e := c.GetData(ctx, "counter", &out); e != nil || rev != 1 || out != 4 {
		t.Fatal(e, rev, out)
	}
	if _, e = c.SaveData(ctx, "counter", 5, 0); e == nil {
		t.Fatal("stale KV revision accepted")
	}
	if path, e := c.GetDataPath(ctx); e != nil || path != "/fixture/plugins/instance/data" {
		t.Fatal(e, path)
	}
	if e = c.PostMessage(ctx, contracts.Message{Title: "hello", Text: "test"}); e != nil || len(m.Messages()) != 1 {
		t.Fatal(e)
	}
	if e = c.DelData(ctx, "counter", 1); e != nil {
		t.Fatal(e)
	}
	var rpc *contracts.Error
	if _, e = c.GetData(ctx, "counter", &out); !errors.As(e, &rpc) || rpc.Code != contracts.NotFound {
		t.Fatal(e)
	}
	m.Revoke("plugin.data.write")
	if _, e = c.SaveData(ctx, "counter", 1, 0); !errors.As(e, &rpc) || rpc.Code != contracts.PermissionDenied {
		t.Fatal(e)
	}
	if _, e = c.SaveData(ctx, strings.Repeat("x", 129), 1, 0); !errors.As(e, &rpc) || rpc.Code != contracts.InvalidParams {
		t.Fatal(e)
	}
	if _, e = c.SaveData(ctx, "counter", strings.Repeat("x", 128<<10), 0); !errors.As(e, &rpc) || rpc.Code != contracts.InvalidParams {
		t.Fatal(e)
	}
	if _, e = c.UpdateConfig(ctx, map[string]int{}, -1); !errors.As(e, &rpc) || rpc.Code != contracts.InvalidParams {
		t.Fatal(e)
	}
	m.InjectOnce("plugin.data.get", contracts.UpstreamLimited, 0)
	if _, e = c.GetData(ctx, "x", &out); !errors.As(e, &rpc) || rpc.Code != contracts.UpstreamLimited {
		t.Fatal(e)
	}
	m.InjectOnce("plugin.config.get", 0, time.Second)
	short, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	if _, e = c.GetConfig(short, &config); !errors.Is(e, context.DeadlineExceeded) {
		t.Fatal(e)
	}
}
