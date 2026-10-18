// Package cache provides a bounded, credential-free platform Host client.
package cache

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Status(ctx context.Context, in b.NamespaceQuery) (b.CacheSummary, error) {
	return b.Call[b.CacheSummary](ctx, c.Host, "cache.status", in, b.Operation{})
}
func (c *Client) List(ctx context.Context, in b.NamespaceQuery) (b.CachePage, error) {
	return b.Call[b.CachePage](ctx, c.Host, "cache.list", in, b.Operation{})
}
func (c *Client) Invalidate(ctx context.Context, in b.CacheInvalidation, op b.Operation) (b.Count, error) {
	return b.Call[b.Count](ctx, c.Host, "cache.invalidate", in, op)
}
