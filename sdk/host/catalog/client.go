// Package catalog exposes safe typed media Host APIs; no playback URL or member credential methods.
package catalog

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) List(ctx context.Context, in b.ServerQuery) (b.Page[b.ServerItem], error) {
	return b.Call[b.Page[b.ServerItem]](ctx, c.Host, "catalog.list", in, b.Operation{})
}
func (c *Client) Item(ctx context.Context, in b.ServerQuery) (b.ServerItem, error) {
	return b.Call[b.ServerItem](ctx, c.Host, "catalog.item", in, b.Operation{})
}
func (c *Client) Home(ctx context.Context, in b.ServerQuery) (b.Page[b.ServerItem], error) {
	return b.Call[b.Page[b.ServerItem]](ctx, c.Host, "catalog.home", in, b.Operation{})
}
func (c *Client) Discover(ctx context.Context, in b.ServerQuery) (b.Page[b.ServerItem], error) {
	return b.Call[b.Page[b.ServerItem]](ctx, c.Host, "catalog.discover", in, b.Operation{})
}
func (c *Client) NextUp(ctx context.Context, in b.ServerQuery) (b.Page[b.ServerItem], error) {
	return b.Call[b.Page[b.ServerItem]](ctx, c.Host, "catalog.nextUp", in, b.Operation{})
}
func (c *Client) UserData(ctx context.Context, in b.ServerQuery) (b.Empty, error) {
	return b.Call[b.Empty](ctx, c.Host, "catalog.userData", in, b.Operation{})
}
