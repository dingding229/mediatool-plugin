package plugintest

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	"github.com/dingding229/mediatool-plugin/sdk/transport"
	"net"
	"sync"
	"time"
)

type Harness struct {
	Peer      *transport.Peer
	Client    *host.Client
	done      chan error
	worker    *transport.Peer
	closeOnce sync.Once
	closeErr  error
}

func Start(ctx context.Context, m contracts.Manifest, f sdk.Factory, h *Host, id sdk.Identity) (*Harness, error) {
	a, b := net.Pipe()
	hp, e := transport.New(a, transport.Options{Prefix: "h"})
	if e != nil {
		return nil, e
	}
	wp, e := transport.New(b, transport.Options{Prefix: "p"})
	if e != nil {
		hp.Close()
		b.Close()
		return nil, e
	}
	if e = hp.Handle("host.call", h.Handle); e != nil {
		hp.Close()
		wp.Close()
		return nil, e
	}
	r, e := sdk.NewRuntime(wp, m, f, id)
	if e != nil {
		hp.Close()
		wp.Close()
		return nil, e
	}
	if e = hp.Start(ctx); e != nil {
		hp.Close()
		wp.Close()
		return nil, e
	}
	harness := &Harness{Peer: hp, Client: host.New(wp), done: make(chan error, 1), worker: wp}
	go func() { harness.done <- r.Serve(ctx) }()
	select {
	case <-wp.Ready():
		return harness, nil
	case <-ctx.Done():
		harness.Close()
		return nil, ctx.Err()
	}
}
func (h *Harness) Initialize(ctx context.Context, q contracts.Initialize) (contracts.Initialized, error) {
	var out contracts.Initialized
	e := h.Peer.Call(ctx, "plugin.initialize", q, &out)
	return out, e
}
func (h *Harness) Init(ctx context.Context, c contracts.Config) error {
	return h.Peer.Call(ctx, "plugin.init", c, nil)
}
func (h *Harness) Call(ctx context.Context, method string, args, out any) error {
	return h.Peer.Call(ctx, method, args, out)
}
func (h *Harness) Close() error {
	h.closeOnce.Do(func() {
		h.Peer.Close()
		h.worker.Close()
		select {
		case <-h.done:
		case <-time.After(12 * time.Second):
			h.closeErr = contracts.NewError(contracts.NeedsAction)
		}
	})
	return h.closeErr
}

// JSON is convenient for synthetic test values and refuses unsupported values.
func JSON(v any) (json.RawMessage, error) { return json.Marshal(v) }
