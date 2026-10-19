package seeding

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func fixture(t *testing.T, execute, files bool) (*Plugin, *plugintest.Host, *int) {
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
	raw, _ := json.Marshal(Config{Active: true, Execute: execute, DeleteFiles: files, MinRatio: 2, MinSeedingHours: 1, BatchSize: 2})
	if e = p.InitPlugin(context.Background(), sdk.Config{Value: raw}); e != nil {
		t.Fatal(e)
	}
	n := new(int)
	task := b.Task{ID: "task", Hash: "hash", Owned: true, Complete: true, Seeding: true, CleanupEligible: true, Ratio: 3, SeedingSeconds: 7200}
	h.OnBusiness("downloads.list", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.Page[b.Task]{Items: []b.Task{task, task, task}}, nil
	})
	h.OnBusiness("downloads.downloader.inspect", func(context.Context, json.RawMessage) (any, *contracts.Error) { return task, nil })
	h.OnBusiness("downloads.downloader.deleteOwned", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		*n++
		var q b.DeleteRequest
		if contracts.Decode(raw, &q) != nil || q.DeleteFiles != files {
			t.Error(string(raw))
		}
		return b.Receipt{State: "retired"}, nil
	})
	return p, h, n
}
func TestDryRunBoundedExplicitExecutionAndReceipt(t *testing.T) {
	for _, execute := range []bool{false, true} {
		p, h, n := fixture(t, execute, false)
		for range 2 {
			if e := p.run(context.Background(), nil); e != nil {
				t.Fatal(e)
			}
		}
		want := 0
		if execute {
			want = 1
		}
		if *n != want {
			t.Fatal(*n, want)
		}
		var s State
		if _, e := h.Client().GetData(context.Background(), "cleanup-state", &s); e != nil || len(s.Items) != 2 {
			t.Fatal(e, s)
		}
	}
}
func TestSafetyRulesAndRevocation(t *testing.T) {
	p, h, n := fixture(t, true, false)
	base := b.Task{ID: "id", Owned: true, Complete: true, Seeding: true, CleanupEligible: true, Ratio: 3, SeedingSeconds: 7200}
	if !p.eligible(base) {
		t.Fatal("valid task rejected")
	}
	for _, mutate := range []func(*b.Task){func(x *b.Task) { x.Owned = false }, func(x *b.Task) { x.Complete = false }, func(x *b.Task) { x.CleanupEligible = false }, func(x *b.Task) { x.Seeding = false }, func(x *b.Task) { x.Ratio = 1 }, func(x *b.Task) { x.SeedingSeconds = 1 }} {
		v := base
		mutate(&v)
		if p.eligible(v) {
			t.Fatal("unsafe task eligible", v)
		}
	}
	p.config.Tag = "keep"
	if p.eligible(base) {
		t.Fatal("tag filter ignored")
	}
	p.config.Tag = ""
	h.Revoke("downloads.delete")
	if e := p.run(context.Background(), nil); e == nil || *n != 0 {
		t.Fatal(e, *n)
	}
}
func TestUnknownIntentNeverReplayed(t *testing.T) {
	p, h, n := fixture(t, true, true)
	h.OnBusiness("downloads.downloader.deleteOwned", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.Receipt{State: "unknown"}, nil
	})
	for range 2 {
		if e := p.run(context.Background(), nil); e == nil {
			t.Fatal("unknown accepted")
		}
	}
	if *n != 1 {
		t.Fatal(*n)
	}
}
func TestDefaultInactiveAndFormNumericStrings(t *testing.T) {
	p, h, _ := fixture(t, false, false)
	p.config.Active = false
	if e := p.run(context.Background(), nil); e != nil || len(h.Calls()) != 0 {
		t.Fatal(e, h.Calls())
	}
	if e := p.InitPlugin(context.Background(), sdk.Config{Value: json.RawMessage(`{"minRatio":"2","minSeedingHours":"168","batchSize":"10"}`)}); e != nil {
		t.Fatal(e)
	}
	for _, raw := range []string{`{"batchSize":21}`, `{"batchSize":1.2}`, `{"minRatio":"NaN"}`, `{"minSeedingHours":"NaN"}`} {
		if e := p.InitPlugin(context.Background(), sdk.Config{Value: json.RawMessage(raw)}); e == nil {
			t.Fatal(raw)
		}
	}
}
