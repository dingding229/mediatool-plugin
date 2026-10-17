package events

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"testing"
)

func TestChainPriorityAndFreeze(t *testing.T) {
	r := New()
	if e := r.OnChainPriority(contracts.ChainNameRecognize, 50, func(context.Context, ChainEnvelope) (json.RawMessage, error) {
		return json.RawMessage(`{"continue":false}`), nil
	}); e != nil {
		t.Fatal(e)
	}
	if r.Descriptors()[0].Priority != 50 {
		t.Fatal("priority lost")
	}
	r.Freeze()
	if r.OnChain(contracts.ChainNameRecognize, nil) == nil {
		t.Fatal("late handler")
	}
	out, e := r.DispatchChain(context.Background(), ChainEnvelope{Type: contracts.ChainNameRecognize, SchemaVersion: 1, Payload: json.RawMessage(`{}`)})
	if e != nil || string(out) != `{"continue":false}` {
		t.Fatal(string(out), e)
	}
}
