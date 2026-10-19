package plugintest

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
)

// Client is an in-process typed mock; use Start for real transport coverage.
func (h *Host) Client() *host.Client { return host.New(mockCaller{h}) }

type mockCaller struct{ h *Host }

func (c mockCaller) Call(ctx context.Context, method string, in, out any) error {
	if method != "host.call" {
		return contracts.NewError(contracts.MethodNotFound)
	}
	raw, e := json.Marshal(in)
	if e != nil {
		return e
	}
	v, rpc := c.h.Handle(ctx, raw)
	if rpc != nil {
		return rpc
	}
	if out == nil {
		return nil
	}
	raw, e = json.Marshal(v)
	if e != nil {
		return e
	}
	return json.Unmarshal(raw, out)
}
