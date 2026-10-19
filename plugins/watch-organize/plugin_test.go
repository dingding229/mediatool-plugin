package watch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/files"
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
		return b.FilePlan{ID: "plan", Revision: 3, State: "preview", Digest: strings.Repeat("a", 64)}, nil
	})
	h.OnBusiness("fs.media.apply", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.FilePlan{ID: "plan", State: "succeeded"}, nil
	})
	h.OnBusiness("fs.media.get", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.FilePlan{ID: "plan", State: "succeeded"}, nil
	})
	p.host = host.New(revisionCaller{h: h})
	return p, h, f, n
}
func TestQuietCoalesceRestartAndDuplicate(t *testing.T) {
	p, _, f, n := fixture(t)
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
	restarted := New(sdk.Services{Host: p.host}).(*Plugin)
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

func TestMissingMtimeIsExplicitUnsupportedNotInventedStability(t *testing.T) {
	p, _, f, n := fixture(t)
	f.Modified = 0
	err := p.run(context.Background(), nil)
	var rpc *contracts.Error
	if !errors.As(err, &rpc) || rpc.Code != contracts.Unsupported || *n != 0 {
		t.Fatal(err, *n)
	}
}

// The core apply contract requires the preview root revision as well as digest.
// Keep the mock honest so a nil revision cannot appear successful again.
type revisionCaller struct{ h *plugintest.Host }

func (c revisionCaller) Call(ctx context.Context, method string, in, out any) error {
	q, ok := in.(contracts.HostCall)
	if !ok {
		return contracts.NewError(contracts.InvalidParams)
	}
	if q.Method == "fs.media.apply" && (q.ExpectedRevision == nil || *q.ExpectedRevision != 3) {
		return contracts.NewError(contracts.Conflict)
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

func TestBoundedRecursiveLeafQuietDepthAndLinkRefusal(t *testing.T) {
	p, h, _, n := fixture(t)
	calls := []string{}
	h.OnBusiness("fs.media.list", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		var q b.FileQuery
		json.Unmarshal(raw, &q)
		calls = append(calls, q.Path)
		items := []b.MediaFile{}
		switch q.Path {
		case "incoming":
			items = []b.MediaFile{{ID: "show", Name: "Show", Path: "incoming/Show", Directory: true}, {ID: "link", Name: "escape", Path: "incoming/escape", Directory: true, Unsupported: true}}
		case "incoming/Show":
			items = []b.MediaFile{{ID: "season", Name: "Season 01", Path: "incoming/Show/Season 01", Directory: true}}
		case "incoming/Show/Season 01":
			items = []b.MediaFile{{ID: "episode", Name: "Show.S01E01.mkv", Path: "incoming/Show/Season 01/Show.S01E01.mkv", Size: 4, Modified: 1}, {ID: "deep", Name: "deep", Path: "incoming/Show/Season 01/deep", Directory: true}}
		default:
			t.Error("escaped or exceeded depth", q.Path)
		}
		return b.FilePage{Page: b.Page[b.MediaFile]{Items: items}}, nil
	})
	h.OnBusiness("fs.media.preview", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		var in b.FileInput
		json.Unmarshal(raw, &in)
		if in.Source != "incoming/Show/Season 01/Show.S01E01.mkv" || in.ExpectedModified == nil {
			t.Error(in)
		}
		return b.FilePlan{ID: "plan", Revision: 3, State: "preview", Digest: strings.Repeat("a", 64)}, nil
	})
	if e := p.run(context.Background(), nil); e != nil || *n != 0 {
		t.Fatal(e, *n)
	}
	p.now = func() time.Time { return time.Unix(1000031, 0) }
	if e := p.run(context.Background(), nil); e != nil || *n != 1 {
		t.Fatal(e, *n)
	}
	if len(calls) != 6 {
		t.Fatal(calls)
	}
	p.config.Depth = 1
	leaf, e := p.collect(context.Background(), files.New(p.host))
	if e != nil || len(leaf) != 0 {
		t.Fatal(leaf, e)
	}
}
func TestRecursiveTotalAndAliasBound(t *testing.T) {
	p, h, _, _ := fixture(t)
	h.OnBusiness("fs.media.list", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		items := []b.MediaFile{}
		for i := range 201 {
			items = append(items, b.MediaFile{ID: fmt.Sprint(i), Path: "incoming/file"})
		}
		return b.FilePage{Page: b.Page[b.MediaFile]{Items: items}}, nil
	})
	if e := p.run(context.Background(), nil); e == nil {
		t.Fatal("unbounded recursive scan")
	}
	for _, v := range []string{`{"depth":0}`, `{"depth":6}`, `{"depth":2.5}`} {
		if e := p.InitPlugin(context.Background(), sdk.Config{Value: json.RawMessage(v)}); e == nil {
			t.Fatal(v)
		}
	}
}
