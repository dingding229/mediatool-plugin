package watch

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"strings"
	"testing"
	"time"
)

func fixture(t *testing.T) (*Plugin, *plugintest.Host, *b.MediaFile, *int) {
	t.Helper()
	m, e := contracts.ParseManifest(ManifestJSON())
	if e != nil {
		t.Fatal(e)
	}
	scopes := []string{}
	for _, s := range m.Permissions {
		scopes = append(scopes, s.Scope)
	}
	h := plugintest.NewHost(scopes)
	p := New(sdk.Services{Host: h.Client()}).(*Plugin)
	p.now = func() time.Time { return time.Unix(1000000, 0) }
	if e = p.InitPlugin(context.Background(), sdk.Config{Value: json.RawMessage(`{"active":true,"rootId":"00000000-0000-0000-0000-000000000001","directory":"incoming","destination":"organized","quietSeconds":"30"}`)}); e != nil {
		t.Fatal(e)
	}
	f := &b.MediaFile{ID: "file-id", Name: "sample.mkv", Path: "incoming/sample.mkv", Size: 4, Modified: 1}
	n := new(int)
	h.OnBusiness("fs.media.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.FilePage{Page: b.Page[b.MediaFile]{Items: []b.MediaFile{*f}}}, nil
	})
	h.OnBusiness("fs.media.preview", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		var in b.FileInput
		if contracts.Decode(raw, &in) != nil || in.Operation != "organize" || in.Name != "" || in.ExpectedModified == nil || *in.ExpectedModified != f.Modified || in.ExpectedFileID != f.ID || in.ExpectedSize == nil {
			t.Error(string(raw))
		}
		return b.FilePlan{ID: "plan", State: "preview", Digest: strings.Repeat("a", 64)}, nil
	})
	h.OnBusiness("fs.media.apply", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.FilePlan{ID: "plan", State: "succeeded"}, nil
	})
	h.OnBusiness("fs.media.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.FilePlan{ID: "plan", State: "succeeded"}, nil
	})
	return p, h, f, n
}
func TestQuietCoalesceRestartAndDuplicate(t *testing.T) {
	p, h, f, n := fixture(t)
	ctx := context.Background()
	if e := p.run(ctx, nil); e != nil || *n != 0 {
		t.Fatal(e, *n)
	}
	p.now = func() time.Time { return time.Unix(1000020, 0) }
	f.Size++
	f.Modified++
	if e := p.run(ctx, nil); e != nil {
		t.Fatal(e)
	}
	p.now = func() time.Time { return time.Unix(1000040, 0) }
	if e := p.run(ctx, nil); e != nil || *n != 0 {
		t.Fatal(e, *n)
	}
	// A fresh worker uses the persisted quiet observation, rather than losing its queue.
	restarted := New(sdk.Services{Host: h.Client()}).(*Plugin)
	restarted.config = p.config
	restarted.running = true
	restarted.now = func() time.Time { return time.Unix(1000050, 0) }
	for range 2 {
		if e := restarted.run(ctx, nil); e != nil {
			t.Fatal(e)
		}
	}
	if *n != 1 {
		t.Fatal(*n)
	}
}
func TestOverflowUnsafePathsRevocationAndUnknown(t *testing.T) {
	p, h, f, n := fixture(t)
	ctx := context.Background()
	h.OnBusiness("fs.media.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.FilePage{Page: b.Page[b.MediaFile]{HasMore: true}}, nil
	})
	if e := p.run(ctx, nil); e == nil || *n != 0 {
		t.Fatal(e)
	}
	p, h, f, n = fixture(t)
	f.Path = "../escape.mkv"
	if e := p.run(ctx, nil); e == nil {
		t.Fatal("escape accepted")
	}
	p, h, f, n = fixture(t)
	f.Unsupported = true
	for range 2 {
		if e := p.run(ctx, nil); e != nil {
			t.Fatal(e)
		}
	}
	if *n != 0 {
		t.Fatal("unsupported file moved")
	}
	p, h, f, n = fixture(t)
	if e := p.run(ctx, nil); e != nil {
		t.Fatal(e)
	}
	p.now = func() time.Time { return time.Unix(1000031, 0) }
	h.Revoke("fs.media:write")
	if e := p.run(ctx, nil); e == nil || *n != 0 {
		t.Fatal(e, *n)
	}
	p, h, f, n = fixture(t)
	p.run(ctx, nil)
	p.now = func() time.Time { return time.Unix(1000031, 0) }
	h.OnBusiness("fs.media.apply", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.FilePlan{State: "unknown"}, nil
	})
	h.OnBusiness("fs.media.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.FilePlan{State: "uncertain"}, nil
	})
	for range 2 {
		if e := p.run(ctx, nil); e == nil {
			t.Fatal("unknown replay")
		}
	}
	if *n != 1 {
		t.Fatal(*n)
	}
}
func TestPathsAndInactive(t *testing.T) {
	p, h, _, _ := fixture(t)
	for _, dir := range []string{"../outside", "/absolute", "incoming/../outside", "a\\b"} {
		raw, _ := json.Marshal(Config{Active: true, RootID: "root", Directory: dir, Destination: "out", QuietSeconds: 30})
		if e := p.InitPlugin(context.Background(), sdk.Config{Value: raw}); e == nil {
			t.Fatal(dir)
		}
	}
	p.config.Active = false
	if e := p.run(context.Background(), nil); e != nil || len(h.Calls()) != 0 {
		t.Fatal(e)
	}
}
