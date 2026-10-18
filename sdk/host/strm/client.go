// Package strm exposes typed, credential-free media/file Host APIs.
package strm

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Entries(ctx context.Context, in b.PageRequest) (b.Page[b.Asset], error) {
	return b.Call[b.Page[b.Asset]](ctx, c.Host, "strm.entries", in, b.Operation{})
}
func (c *Client) Diagnostics(ctx context.Context, in b.ID) (b.RepairPlan, error) {
	return b.Call[b.RepairPlan](ctx, c.Host, "strm.diagnostics", in, b.Operation{})
}
func (c *Client) Preview(ctx context.Context, in b.IDs, op b.Operation) (b.RepairPlan, error) {
	return b.Call[b.RepairPlan](ctx, c.Host, "strm.repair.preview", in, op)
}
func (c *Client) Approve(ctx context.Context, in b.Approval, op b.Operation) (b.RepairPlan, error) {
	return b.Call[b.RepairPlan](ctx, c.Host, "strm.repair.approve", in, op)
}
func (c *Client) Retry(ctx context.Context, in b.Approval, op b.Operation) (b.RepairPlan, error) {
	return b.Call[b.RepairPlan](ctx, c.Host, "strm.repair.retry", in, op)
}
