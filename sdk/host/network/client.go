// Package network provides a bounded, credential-free platform Host client.
package network

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) ProxyStatus(ctx context.Context, in b.Empty) (b.ProxySummary, error) {
	return b.Call[b.ProxySummary](ctx, c.Host, "networkproxy.status", in, b.Operation{})
}
func (c *Client) Request(ctx context.Context, in b.HostRequest, op b.Operation) (b.HTTPResult, error) {
	return b.Call[b.HTTPResult](ctx, c.Host, "network.http", in, op)
}
