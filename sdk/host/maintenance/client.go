// Package maintenance provides a bounded, credential-free platform Host client.
package maintenance

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Health(ctx context.Context, in b.Empty) (b.HealthSummary, error) {
	return b.Call[b.HealthSummary](ctx, c.Host, "maintenance.health", in, b.Operation{})
}
func (c *Client) Alerts(ctx context.Context, in b.PageRequest) (b.Page[b.AlertSummary], error) {
	return b.Call[b.Page[b.AlertSummary]](ctx, c.Host, "maintenance.alerts", in, b.Operation{})
}
