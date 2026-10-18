// Package classification exposes typed, credential-free media/file Host APIs.
package classification

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Get(ctx context.Context, in b.Empty) (b.Classification, error) {
	return b.Call[b.Classification](ctx, c.Host, "classification.get", in, b.Operation{})
}
func (c *Client) Match(ctx context.Context, in b.Classify) (b.Category, error) {
	return b.Call[b.Category](ctx, c.Host, "classification.match", in, b.Operation{})
}
func (c *Client) Render(ctx context.Context, in b.Classify) (string, error) {
	return b.Call[string](ctx, c.Host, "classification.render", in, b.Operation{})
}
func (c *Client) Update(ctx context.Context, in b.Classification, op b.Operation) (b.Classification, error) {
	return b.Call[b.Classification](ctx, c.Host, "classification.update", in, op)
}
