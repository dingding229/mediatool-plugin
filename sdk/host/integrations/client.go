// Package integrations provides a bounded, credential-free platform Host client.
package integrations

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) List(ctx context.Context, in b.PageRequest) (b.Page[b.Connection], error) {
	return b.Call[b.Page[b.Connection]](ctx, c.Host, "integrations.list", in, b.Operation{})
}
func (c *Client) Status(ctx context.Context, in b.ID) (b.Connection, error) {
	return b.Call[b.Connection](ctx, c.Host, "integrations.status", in, b.Operation{})
}
func (c *Client) Test(ctx context.Context, in b.ID, op b.Operation) (b.ConnectionCheck, error) {
	return b.Call[b.ConnectionCheck](ctx, c.Host, "integrations.test", in, op)
}
