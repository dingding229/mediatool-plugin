package sdk_test

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"testing"
)

func TestHelloAcquisitionNoopAndTypedEventReceipts(t *testing.T) {
	d, m, id, q := fixture(t)
	h := start(t, d, m, id)
	ctx := context.Background()
	if _, e := h.Initialize(ctx, q); e != nil {
		t.Fatal(e)
	}
	if e := h.Init(ctx, contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}); e != nil {
		t.Fatal(e)
	}
	a := "00000000-0000-0000-0000-000000000001"
	selection := contracts.ResourceSelectionRequest{RunID: a, Candidates: []contracts.ResourceCandidate{{ID: a, AccountID: a}}}
	raw, _ := json.Marshal(selection)
	var out contracts.ResourceSelectionResult
	if e := h.Call(ctx, "plugin.invoke", contracts.Invocation{Handler: contracts.HandlerID("chain", "resource.selection"), Args: mustJSON(events.ChainEnvelope{Type: contracts.ChainResourceSelection, SchemaVersion: 1, Payload: raw})}, &out); e != nil || contracts.ValidateSelection(selection, out) != nil || len(out.IDs) != 1 {
		t.Fatal(out, e)
	}
	download := contracts.ResourceDownloadRequest{Kind: "managed", ResourceID: a, DownloaderID: a, CleanupAfterImport: true}
	raw, _ = json.Marshal(download)
	var result contracts.ResourceDownloadResult
	if e := h.Call(ctx, "plugin.invoke", contracts.Invocation{Handler: contracts.HandlerID("chain", "resource.download"), Args: mustJSON(events.ChainEnvelope{Type: contracts.ChainResourceDownload, SchemaVersion: 1, Payload: raw})}, &result); e != nil || !result.Allow || !result.CleanupAfterImport {
		t.Fatal(result, e)
	}
	for _, typ := range []contracts.EventType{contracts.EventSiteUpdated, contracts.EventSiteDeleted, contracts.EventSiteRefreshed, contracts.EventSubscribeAdded, contracts.EventSubscribeModified, contracts.EventSubscribeDeleted} {
		for range 2 {
			if e := h.Call(ctx, "plugin.invoke", contracts.Invocation{Handler: contracts.HandlerID("event", string(typ)), Args: mustJSON(events.Envelope{ID: string(typ), Type: typ, SchemaVersion: 1, Payload: json.RawMessage(`{}`)})}, nil); e != nil {
				t.Fatal(e)
			}
		}
	}
}
func mustJSON(v any) json.RawMessage {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return b
}
