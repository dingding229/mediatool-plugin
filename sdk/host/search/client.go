// Package search is a typed, credential-free Host client. Host authorization is always live.
package search

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(c *host.Client) *Client { return &Client{Host: c} }
func (c *Client) Start(ctx context.Context, in b.SearchStart, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "search.start", in, op)
}
func (c *Client) List(ctx context.Context, in b.PageRequest) (b.Page[b.SearchRun], error) {
	return b.Call[b.Page[b.SearchRun]](ctx, c.Host, "search.list", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.Selection) (b.SearchRun, error) {
	return b.Call[b.SearchRun](ctx, c.Host, "search.get", in, b.Operation{})
}
func (c *Client) Resources(ctx context.Context, in b.Selection) (b.SearchRun, error) {
	return b.Call[b.SearchRun](ctx, c.Host, "search.resources", in, b.Operation{})
}
