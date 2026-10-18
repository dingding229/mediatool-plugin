// Package system provides a bounded, credential-free platform Host client.
package system

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Info(ctx context.Context, in b.Empty) (b.SystemInfo, error) {
	return b.Call[b.SystemInfo](ctx, c.Host, "system.info", in, b.Operation{})
}
func (c *Client) Health(ctx context.Context, in b.Empty) (b.HealthSummary, error) {
	return b.Call[b.HealthSummary](ctx, c.Host, "system.health", in, b.Operation{})
}
func (c *Client) Modules(ctx context.Context, in b.Empty) ([]b.SystemModule, error) {
	return b.Call[[]b.SystemModule](ctx, c.Host, "system.modules", in, b.Operation{})
}
