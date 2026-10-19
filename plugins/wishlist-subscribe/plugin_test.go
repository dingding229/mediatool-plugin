package wishlist

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"strconv"
	"testing"
	"time"
)

func setup(t *testing.T, create bool) (*Plugin, *plugintest.Host, *int, *int) {
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
	raw, _ := json.Marshal(Config{Active: true, Source: "tmdb-list", ID: "42", AccountIDs: "00000000-0000-0000-0000-000000000001", Create: create, IntervalMinutes: 30, PerRun: 10})
	if e = p.InitPlugin(context.Background(), sdk.Config{Value: raw}); e != nil {
		t.Fatal(e)
	}
	fetch, n := new(int), new(int)
	h.OnBusiness("metadata.wishlist", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*fetch++
		return b.WishlistPage{State: "ok", Items: []b.WishlistEntry{{SourceID: "tmdb:movie:1", TMDBID: 1, MediaType: "movie", Title: "Synthetic", Year: "2026"}}}, nil
	})
	h.OnBusiness("subscriptions.ensure", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		*n++
		var in b.SubscriptionInput
		if contracts.Decode(raw, &in) != nil || in.Enabled || in.Source != "search" {
			t.Error(string(raw))
		}
		return b.Receipt{ID: "subscription", State: "ensured"}, nil
	})
	return p, h, fetch, n
}
func TestDefaultPreviewExplicitCreateAndDedup(t *testing.T) {
	for _, create := range []bool{false, true} {
		p, h, fetch, n := setup(t, create)
		ctx := context.Background()
		if e := p.run(ctx, nil); e != nil {
			t.Fatal(e)
		}
		if e := p.run(ctx, nil); e != nil || *fetch != 1 {
			t.Fatal(e, *fetch)
		}
		p.now = func() time.Time { return time.Unix(1001801, 0) }
		if e := p.run(ctx, nil); e != nil {
			t.Fatal(e)
		}
		want := 0
		if create {
			want = 1
		}
		if *n != want {
			t.Fatal("effect repeated", *n)
		}
		var s State
		if _, e := h.Client().GetData(ctx, "wishlist-state", &s); e != nil {
			t.Fatal(e)
		}
		if !create && len(s.Pending) != 1 {
			t.Fatal(s)
		}
	}
}
func TestAmbiguousManualChoiceAndInjectedIdentity(t *testing.T) {
	p, h, _, n := setup(t, true)
	h.OnBusiness("metadata.wishlist", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.WishlistPage{State: "ok", Items: []b.WishlistEntry{{SourceID: "douban:1", Title: "Same", Year: "2026"}}}, nil
	})
	h.OnBusiness("metadata.search", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
		var q b.MetadataQuery
		json.Unmarshal(raw, &q)
		if q.MediaType == "series" {
			return []b.Media{}, nil
		}
		return []b.Media{{ID: 1, MediaType: "movie", Title: "Same", Year: "2026"}, {ID: 2, MediaType: "movie", Title: "Same", Year: "2026"}}, nil
	})
	ctx := context.Background()
	if e := p.run(ctx, nil); e != nil || *n != 0 {
		t.Fatal(e, *n)
	}
	for _, id := range []int{999, 2} {
		q := sdk.ActionContext{Data: map[string]json.RawMessage{"sourceId": json.RawMessage(`"douban:1"`), "tmdbId": json.RawMessage(strconv.Itoa(id)), "mediaType": json.RawMessage(`"movie"`)}}
		_, _, e := p.resolve(ctx, q)
		if id == 999 && e == nil || id == 2 && e != nil {
			t.Fatal(id, e)
		}
	}
	if *n != 1 {
		t.Fatal(*n)
	}
}
func TestBackoffCapacityRevocationAndInactive(t *testing.T) {
	p, h, fetch, n := setup(t, true)
	h.OnBusiness("metadata.wishlist", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*fetch++
		return b.WishlistPage{State: "limited"}, nil
	})
	for range 2 {
		if e := p.run(context.Background(), nil); e != nil {
			t.Fatal(e)
		}
	}
	if *fetch != 1 || *n != 0 {
		t.Fatal(*fetch, *n)
	}
	p, h, _, n = setup(t, true)
	h.Revoke("subscriptions.write")
	if e := p.run(context.Background(), nil); e != nil {
		t.Fatal(e)
	}
	var s State
	h.Client().GetData(context.Background(), "wishlist-state", &s)
	if s.Status != "ensure_unconfirmed" || *n != 0 {
		t.Fatal(s, *n)
	}
	p, h, _, _ = setup(t, false)
	p.config.Active = false
	if e := p.run(context.Background(), nil); e != nil || len(h.Calls()) != 0 {
		t.Fatal(e)
	}
}
func TestLimitsAndUnknownAreNotSuccess(t *testing.T) {
	p, h, _, _ := setup(t, true)
	h.OnBusiness("subscriptions.ensure", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		return b.Receipt{State: "unknown"}, nil
	})
	if e := p.run(context.Background(), nil); e != nil {
		t.Fatal(e)
	}
	var s State
	h.Client().GetData(context.Background(), "wishlist-state", &s)
	if s.Ensured != 0 || s.Status != "ensure_unconfirmed" {
		t.Fatal(s)
	}
	for _, raw := range []string{`{"intervalMinutes":1}`, `{"perRun":201}`, `{"source":"external-url"}`} {
		if e := p.InitPlugin(context.Background(), sdk.Config{Value: json.RawMessage(raw)}); e == nil {
			t.Fatal(raw)
		}
	}
}

func TestPendingUIStaysBounded(t *testing.T) {
	p, h, _, _ := setup(t, false)
	s := State{}
	for i := range 200 {
		s.Pending = append(s.Pending, Pending{Entry: b.WishlistEntry{SourceID: strconv.Itoa(i), Title: "Synthetic"}, Candidates: []b.Media{{ID: 1, MediaType: "movie", Title: "Synthetic"}}})
	}
	h.Client().SaveData(context.Background(), "wishlist-state", s, 0)
	page, e := p.GetPage(context.Background(), sdk.PageRequest{})
	if e != nil || len(page.Content) > 12 {
		t.Fatal(e, len(page.Content))
	}
	for _, v := range page.Content {
		if len(v.Props["text"]) > 8192 {
			t.Fatal("unrenderable pending page")
		}
	}
}
