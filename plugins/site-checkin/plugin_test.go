package checkin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
	"time"
)

func setup(t *testing.T, active bool) (*Plugin, *plugintest.Host) {
	t.Helper()
	m, e := contracts.ParseManifest(ManifestJSON())
	if e != nil {
		t.Fatal(e)
	}
	scopes := []string{}
	for _, v := range m.Permissions {
		scopes = append(scopes, v.Scope)
	}
	h := plugintest.NewHost(scopes)
	p := New(sdk.Services{Host: h.Client()}).(*Plugin)
	p.now = func() time.Time { return time.Date(2026, 10, 19, 8, 0, 0, 0, time.UTC) }
	v, _ := json.Marshal(Config{Active: active, AccountID: "00000000-0000-0000-0000-000000000001", Origin: "https://synthetic.invalid", Path: "/attendance.php", Method: "GET", SuccessMarker: "signed", Cron: "0 8 * * *"})
	if e = p.InitPlugin(context.Background(), sdk.Config{Value: v}); e != nil {
		t.Fatal(e)
	}
	return p, h
}
func register(t *testing.T, h *plugintest.Host, method string, fn func(context.Context, json.RawMessage) (any, *contracts.Error)) {
	t.Helper()
	if e := h.OnBusiness(method, fn); e != nil {
		t.Fatal(e)
	}
}
func TestCheckinEvidenceChallengeUnknownAndRestart(t *testing.T) {
	for _, state := range []string{"matched", "needs_action", "unknown", "transport"} {
		t.Run(state, func(t *testing.T) {
			p, h := setup(t, true)
			requests, notices := 0, 0
			register(t, h, "sites.checkin", func(_ context.Context, raw json.RawMessage) (any, *contracts.Error) {
				requests++
				var q b.SiteRequest
				if contracts.Decode(raw, &q) != nil || q.Raw || q.SuccessMarker != "signed" {
					t.Fatal(string(raw))
				}
				if state == "transport" {
					return nil, contracts.NewError(contracts.Unavailable)
				}
				return b.SiteResponse{Status: 200, State: state, Matched: state == "matched"}, nil
			})
			register(t, h, "notifications.publish", func(context.Context, json.RawMessage) (any, *contracts.Error) {
				notices++
				return b.Receipt{State: "queued"}, nil
			})
			err := p.checkin(context.Background(), nil)
			var rpc *contracts.Error
			if state == "matched" && err != nil {
				t.Fatal(err)
			}
			if state != "matched" && (!errors.As(err, &rpc) || rpc.Code != contracts.NeedsAction) {
				t.Fatal(err)
			}
			restarted := New(sdk.Services{Host: h.Client()}).(*Plugin)
			restarted.now = p.now
			raw, _ := json.Marshal(p.config)
			if e := restarted.InitPlugin(context.Background(), sdk.Config{Value: raw}); e != nil {
				t.Fatal(e)
			}
			_ = restarted.checkin(context.Background(), nil)
			if requests != 1 || notices != 1 {
				t.Fatal(requests, notices)
			}
			var got Result
			if _, e := h.Client().GetData(context.Background(), "last-checkin", &got); e != nil {
				t.Fatal(e)
			}
			if got.Day != "20261019" {
				t.Fatal(got)
			}
		})
	}
}
func TestInactiveZeroIOPrecheckAndRevokedScope(t *testing.T) {
	p, h := setup(t, false)
	if e := p.checkin(context.Background(), nil); e != nil || len(h.Calls()) != 0 {
		t.Fatal(e, h.Calls())
	}
	p, h = setup(t, true)
	p.config.PrecheckPath = "/status"
	n := 0
	register(t, h, "sites.request", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		n++
		return b.SiteResponse{Status: 200, State: "matched", Matched: true}, nil
	})
	register(t, h, "notifications.publish", func(context.Context, json.RawMessage) (any, *contracts.Error) { return b.Receipt{State: "queued"}, nil })
	if e := p.checkin(context.Background(), nil); e != nil || n != 1 {
		t.Fatal(e, n)
	}
	p, h = setup(t, true)
	h.Revoke("sites.request")
	n = 0
	register(t, h, "sites.checkin", func(context.Context, json.RawMessage) (any, *contracts.Error) { n++; return b.SiteResponse{}, nil })
	register(t, h, "notifications.publish", func(context.Context, json.RawMessage) (any, *contracts.Error) { return b.Receipt{State: "queued"}, nil })
	if e := p.checkin(context.Background(), nil); e == nil || n != 0 {
		t.Fatal(e, n)
	}
}
func TestConfigRejectsCredentialOriginsAndEscape(t *testing.T) {
	for _, origin := range []string{"https://user:secret@synthetic.invalid", "https://synthetic.invalid/?token=secret"} {
		p, _ := setup(t, false)
		v, _ := json.Marshal(Config{Active: true, AccountID: "00000000-0000-0000-0000-000000000001", Origin: origin, Path: "/", Method: "GET", SuccessMarker: "ok"})
		if p.InitPlugin(context.Background(), sdk.Config{Value: v}) == nil {
			t.Fatal(origin)
		}
	}
}
