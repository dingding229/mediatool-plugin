package events

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"sync/atomic"
	"testing"
)

func TestTypedEventsFreezeAndKnownPairs(t *testing.T) {
	r := New()
	var count atomic.Int32
	handler := Typed(func(ctx context.Context, p struct {
		AssetID string `json:"assetId"`
	}) error {
		if p.AssetID != "fixture" {
			t.Fatal("wrong payload")
		}
		count.Add(1)
		return nil
	})
	if e := r.On(contracts.EventTransferComplete, handler); e != nil {
		t.Fatal(e)
	}
	if e := r.On(contracts.EventTransferComplete, handler); e == nil {
		t.Fatal("duplicate handler")
	}
	if e := r.On(contracts.EventAudioTransferComplete, handler); e == nil {
		t.Fatal("N/A event accepted")
	}
	if e := r.OnChain(contracts.ChainNameRecognize, func(context.Context, ChainEnvelope) (json.RawMessage, error) {
		return json.RawMessage(`{"continue":true}`), nil
	}); e != nil {
		t.Fatal(e)
	}
	r.Freeze()
	if e := r.On(contracts.EventDownloadAdded, handler); e == nil {
		t.Fatal("registration after init")
	}
	good := Envelope{ID: "1", Type: contracts.EventTransferComplete, SchemaVersion: 1, Payload: json.RawMessage(`{"assetId":"fixture"}`)}
	if e := r.Dispatch(context.Background(), good); e != nil || count.Load() != 1 {
		t.Fatal(e)
	}
	good.Payload = json.RawMessage(`{"assetId":"fixture","unknown":true}`)
	if e := r.Dispatch(context.Background(), good); e == nil {
		t.Fatal("unknown payload accepted")
	}
	good.SchemaVersion = 2
	if e := r.Dispatch(context.Background(), good); e == nil {
		t.Fatal("wrong schema")
	}
	if len(r.Descriptors()) != 2 {
		t.Fatal("descriptors")
	}
	if _, e := r.DispatchChain(context.Background(), ChainEnvelope{Type: contracts.ChainNameRecognize, SchemaVersion: 1, Payload: json.RawMessage(`{}`)}); e != nil {
		t.Fatal(e)
	}
}
