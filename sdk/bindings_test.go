package sdk_test

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"os"
	"testing"
)

func TestBoundDescriptorsStableAndInvokeGolden(t *testing.T) {
	d, m, id, q := fixture(t)
	h := start(t, d, m, id)
	ctx := context.Background()
	if _, e := h.Initialize(ctx, q); e != nil {
		t.Fatal(e)
	}
	if e := h.Init(ctx, contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}); e != nil {
		t.Fatal(e)
	}
	var snapshot contracts.Snapshot
	if e := h.Call(ctx, "plugin.describe", nil, &snapshot); e != nil {
		t.Fatal(e)
	}
	if e := contracts.ValidateSnapshot(d.Manifest, q.Capabilities, snapshot); e != nil {
		t.Fatal(e)
	}
	for _, v := range snapshot.Descriptors {
		var args any = map[string]any{}
		switch v.Kind {
		case "api":
			args = contracts.APIRequest{Method: "GET", Path: v.Path}
		case "action":
			args = contracts.ActionContext{Data: map[string]json.RawMessage{}}
		case "chain":
			args = events.ChainEnvelope{Type: contracts.ChainNameRecognize, SchemaVersion: 1, Payload: json.RawMessage(`{"name":"Demo.2025.mkv","hint":{"title":"Demo","year":"2025","mediaType":"movie","season":0,"episodes":[]},"continue":true}`)}
		case "module":
			args = json.RawMessage(`{"name":"Demo.2025.mkv","hint":{"title":"Demo","year":"2025","mediaType":"movie","season":0,"episodes":[]},"continue":true}`)
		case "event":
			args = events.Envelope{ID: "stable-event", Type: contracts.EventType(v.ID), SchemaVersion: 1, Payload: json.RawMessage(`{"operationId":"synthetic"}`)}
		}
		b, e := json.Marshal(args)
		if e != nil {
			t.Fatal(e)
		}
		var output json.RawMessage
		if e = h.Call(ctx, "plugin.invoke", contracts.Invocation{Handler: v.Handler, Args: b}, &output); e != nil {
			t.Fatal(v.Kind, e)
		}
	}
	// Event duplicate uses the same idempotency key; it neither increments nor fails.
	for _, v := range snapshot.Descriptors {
		if v.Kind == "event" {
			b, _ := json.Marshal(events.Envelope{ID: "stable-event", Type: contracts.EventType(v.ID), SchemaVersion: 1, Payload: json.RawMessage(`{}`)})
			if e := h.Call(ctx, "plugin.invoke", contracts.Invocation{Handler: v.Handler, Args: b}, nil); e != nil {
				t.Fatal(e)
			}
		}
	}
	// Reviewed golden is metadata only; no function pointer/address crosses stdio.
	b, e := json.MarshalIndent(snapshot, "", "  ")
	if e != nil {
		t.Fatal(e)
	}
	b = append(b, '\n')
	want, e := os.ReadFile("testdata/hello-descriptors.json")
	if e != nil {
		t.Fatalf("descriptor fixture missing; actual:\n%s", b)
	}
	if string(want) != string(b) {
		t.Fatalf("descriptor golden mismatch; actual:\n%s", b)
	}
	if e = h.Call(ctx, "plugin.close", nil, nil); e != nil {
		t.Fatal(e)
	}
	rpcCode(t, h.Call(ctx, "plugin.invoke", contracts.Invocation{Handler: snapshot.Descriptors[0].Handler, Args: json.RawMessage(`{}`)}, nil), contracts.NeedsAction)
}
func TestDescriptorValidationRejectsConflictsAndScopeExpansion(t *testing.T) {
	d, _, _, q := fixture(t)
	x := contracts.Descriptor{Kind: "service", ID: "tick", Handler: contracts.HandlerID("service", "tick"), Trigger: "interval", IntervalMS: 1000, RequiredScopes: []string{"services.register"}}
	for _, list := range [][]contracts.Descriptor{{x, x}, {func() contracts.Descriptor { v := x; v.RequiredScopes = []string{"system.exec"}; return v }()}, {func() contracts.Descriptor { v := x; v.IntervalMS = 0; return v }()}} {
		if contracts.ValidateSnapshot(d.Manifest, q.Capabilities, contracts.Snapshot{Version: 2, Descriptors: list}) == nil {
			t.Fatal("invalid descriptor accepted")
		}
	}
	for _, p := range []string{"/../x", "//x", "/%2fadmin", "/foo/{id}", "/foo?token=x"} {
		if contracts.ValidAPIPath(p) {
			t.Fatal(p)
		}
	}
}
