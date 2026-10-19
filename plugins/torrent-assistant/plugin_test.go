package assistant

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
	"time"
)

func setup(t *testing.T, write, recheck bool) (*Plugin, *plugintest.Host, *b.Task, *int, *int) {
	t.Helper()
	m, e := contracts.ParseManifest(ManifestJSON())
	if e != nil {
		t.Fatal(e)
	}
	scopes := []string{}
	for _, x := range m.Permissions {
		scopes = append(scopes, x.Scope)
	}
	h := plugintest.NewHost(scopes)
	p := New(sdk.Services{Host: h.Client()}).(*Plugin)
	p.now = func() time.Time { return time.Unix(1000000, 0) }
	raw, _ := json.Marshal(Config{Active: true, Write: write, Recheck: recheck, Tag: "verified", Batch: 5, StallMinutes: 5})
	if e = p.InitPlugin(context.Background(), sdk.Config{Value: raw}); e != nil {
		t.Fatal(e)
	}
	n, r := new(int), new(int)
	progress := .5
	task := &b.Task{ID: "id", Hash: "hash", Owned: true, State: "error", Progress: &progress}
	h.OnBusiness("downloads.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.Page[b.Task]{Items: []b.Task{*task}}, nil
	})
	h.OnBusiness("downloads.downloader.inspect", func(context.Context, json.RawMessage) (any, *contracts.Error) { return *task, nil })
	h.OnBusiness("downloads.downloader.configure", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.Receipt{State: "configured"}, nil
	})
	h.OnBusiness("downloads.recheck", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*r++
		return b.Receipt{State: "requested"}, nil
	})
	h.OnBusiness("notifications.publish", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.Receipt{State: b.NotificationPending}, nil
	})
	return p, h, task, n, r
}
func TestSuggestDefaultExplicitWritesRecheckAndIdempotence(t *testing.T) {
	for _, write := range []bool{false, true} {
		p, _, _, n, r := setup(t, write, true)
		for range 2 {
			if e := p.run(context.Background(), nil); e != nil {
				t.Fatal(e)
			}
		}
		want := 0
		if write {
			want = 1
		}
		if *n != want || *r != want {
			t.Fatal(*n, *r, want)
		}
	}
}
func TestForeignRevocationAndUnknown(t *testing.T) {
	p, _, task, n, _ := setup(t, true, true)
	task.Owned = false
	if e := p.run(context.Background(), nil); e != nil || *n != 0 {
		t.Fatal(e, *n)
	}
	p, h, _, n, _ := setup(t, true, true)
	h.Revoke("downloads.control")
	if e := p.run(context.Background(), nil); e == nil || *n != 0 {
		t.Fatal(e, *n)
	}
	p, h, _, n, _ = setup(t, true, true)
	h.OnBusiness("downloads.downloader.configure", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.Receipt{State: "unknown"}, nil
	})
	for range 2 {
		if e := p.run(context.Background(), nil); e == nil {
			t.Fatal("unknown replay")
		}
	}
	if *n != 1 {
		t.Fatal(*n)
	}
}
func TestStallProgressAndNumericBounds(t *testing.T) {
	p, _, task, _, r := setup(t, true, true)
	task.State = "stalledDL"
	if e := p.run(context.Background(), nil); e != nil || *r != 0 {
		t.Fatal(e, *r)
	}
	p.now = func() time.Time { return time.Unix(1000301, 0) }
	if e := p.run(context.Background(), nil); e != nil || *r != 1 {
		t.Fatal(e, *r)
	}
	for _, raw := range []string{`{"batch":30}`, `{"batch":1.5}`, `{"ratio":"NaN"}`, `{"downloadKiB":102401}`, `{"stallMinutes":1}`} {
		if e := p.InitPlugin(context.Background(), sdk.Config{Value: json.RawMessage(raw)}); e == nil {
			t.Fatal(raw)
		}
	}
}
func TestInactiveNoConnections(t *testing.T) {
	p, h, _, _, _ := setup(t, false, false)
	p.config.Active = false
	if e := p.run(context.Background(), nil); e != nil || len(h.Calls()) != 0 {
		t.Fatal(e)
	}
}

func TestStatusUIIsBounded(t *testing.T) {
	p, h, _, _, _ := setup(t, false, false)
	s := State{}
	for range 200 {
		s.Items = append(s.Items, Sample{ID: "id", Hash: "hash", State: "preview"})
	}
	h.Client().SaveData(context.Background(), "assistant-state", s, 0)
	page, e := p.GetPage(context.Background(), sdk.PageRequest{})
	if e != nil || len(page.Content) != 21 {
		t.Fatal(e, len(page.Content))
	}
	for _, v := range page.Content {
		if len(v.Props["text"]) > 8192 {
			t.Fatal("unrenderable status")
		}
	}
}
