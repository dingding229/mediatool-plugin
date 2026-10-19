package refresh

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
	"time"
)

func setup(t *testing.T, active bool) (*Plugin, *plugintest.Host, *int, *int) {
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
	p := New(sdk.Services{Host: h.Client(), Events: events.New()}).(*Plugin)
	p.now = func() time.Time { return time.Unix(2000000, 0) }
	raw, _ := json.Marshal(Config{active, "builtin", 60})
	if e = p.InitPlugin(context.Background(), sdk.Config{Value: raw}); e != nil {
		t.Fatal(e)
	}
	n, r := new(int), new(int)
	if e = h.OnBusiness("notifications.publish", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*n++
		return b.Receipt{State: "queued"}, nil
	}); e != nil {
		t.Fatal(e)
	}
	if e = h.OnBusiness("mediaserver.refresh", func(context.Context, json.RawMessage) (any, *contracts.Error) {
		*r++
		return b.Receipt{State: "requested"}, nil
	}); e != nil {
		t.Fatal(e)
	}
	return p, h, n, r
}
func event(n int, asset string) events.Envelope {
	raw, _ := json.Marshal(contracts.TransferCompleted{AssetID: asset, State: "succeeded"})
	return events.Envelope{ID: fmt.Sprintf("00000000-0000-0000-0000-%012d", n), Type: contracts.EventTransferComplete, SchemaVersion: 1, Payload: raw}
}

const asset = "00000000-0000-0000-0000-000000000999"

func TestDuplicateEventsMergeThrottleAndCheckpointRecovery(t *testing.T) {
	p, h, n, r := setup(t, true)
	ctx := context.Background()
	for _, e := range []events.Envelope{event(1, asset), event(1, asset), event(2, asset)} {
		if err := p.onTransfer(ctx, e); err != nil {
			t.Fatal(err)
		}
	}
	if *n != 2 {
		t.Fatal(*n)
	}
	var before State
	rev, e := h.Client().GetData(ctx, "refresh-state", &before)
	if e != nil || len(before.Pending) != 1 {
		t.Fatal(e, before)
	}
	if e = p.flush(ctx, nil); e != nil {
		t.Fatal(e)
	}
	if e = p.flush(ctx, nil); e != nil || *r != 1 {
		t.Fatal(e, *r)
	}
	// Simulate crash after the effect receipt, before the plugin checkpoint.
	var current State
	rev, e = h.Client().GetData(ctx, "refresh-state", &current)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = h.Client().SaveData(ctx, "refresh-state", before, rev); e != nil {
		t.Fatal(e)
	}
	if e = p.flush(ctx, nil); e != nil || *r != 1 {
		t.Fatal("effect replayed", e, *r)
	}
	next := event(3, "00000000-0000-0000-0000-000000000998")
	if e = p.onTransfer(ctx, next); e != nil {
		t.Fatal(e)
	}
	if e = p.flush(ctx, nil); e != nil || *r != 1 {
		t.Fatal("throttle", e, *r)
	}
	p.now = func() time.Time { return time.Unix(2000060, 0) }
	if e = p.flush(ctx, nil); e != nil || *r != 2 {
		t.Fatal(e, *r)
	}
}
func TestUnknownRevocationCapacityAndInactive(t *testing.T) {
	p, h, n, r := setup(t, false)
	if e := p.onTransfer(context.Background(), event(1, asset)); e != nil || len(h.Calls()) != 0 {
		t.Fatal(e)
	}
	p, h, n, r = setup(t, true)
	if e := p.onTransfer(context.Background(), event(1, "")); e != nil || *n != 0 {
		t.Fatal(e, *n)
	}
	if e := p.onTransfer(context.Background(), event(2, asset)); e != nil {
		t.Fatal(e)
	}
	h.Revoke("mediaserver.refresh")
	if e := p.flush(context.Background(), nil); e == nil || *r != 0 {
		t.Fatal(e, *r)
	}
	p, h, _, _ = setup(t, true)
	s := State{}
	for n := range 200 {
		s.Events = append(s.Events, Notice{ID: fmt.Sprint(n), Notified: true})
	}
	if _, e := h.Client().SaveData(context.Background(), "refresh-state", s, 0); e != nil {
		t.Fatal(e)
	}
	if e := p.onTransfer(context.Background(), event(201, asset)); e == nil {
		t.Fatal("unbounded ledger")
	}
}
