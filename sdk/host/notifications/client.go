// Package notifications provides a bounded, credential-free platform Host client.
package notifications

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Send(ctx context.Context, in b.Notice, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "notifications.publish", in, op)
}
func (c *Client) Deliveries(ctx context.Context, in b.PageRequest) (b.Page[b.DeliverySummary], error) {
	return b.Call[b.Page[b.DeliverySummary]](ctx, c.Host, "notifications.deliveries", in, b.Operation{})
}
func (c *Client) Retry(ctx context.Context, in b.DeliveryRetry, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "notifications.retry", in, op)
}
