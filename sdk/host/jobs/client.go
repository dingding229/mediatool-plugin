// Package jobs provides a bounded, credential-free platform Host client.
package jobs

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) List(ctx context.Context, in b.JobQuery) (b.Page[b.OwnedJob], error) {
	return b.Call[b.Page[b.OwnedJob]](ctx, c.Host, "jobs.list", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.ID) (b.OwnedJob, error) {
	return b.Call[b.OwnedJob](ctx, c.Host, "jobs.get", in, b.Operation{})
}
func (c *Client) Cancel(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "jobs.cancel", in, op)
}
