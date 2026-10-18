// Package mediaserver exposes safe typed media Host APIs; no playback URL or member credential methods.
package mediaserver

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Providers(ctx context.Context, in b.Empty) ([]b.ServerProvider, error) {
	return b.Call[[]b.ServerProvider](ctx, c.Host, "mediaserver.providers", in, b.Operation{})
}
func (c *Client) Libraries(ctx context.Context, in b.ServerQuery) (b.Page[b.ServerLibrary], error) {
	return b.Call[b.Page[b.ServerLibrary]](ctx, c.Host, "mediaserver.libraries", in, b.Operation{})
}
func (c *Client) Items(ctx context.Context, in b.ServerQuery) (b.Page[b.ServerItem], error) {
	return b.Call[b.Page[b.ServerItem]](ctx, c.Host, "mediaserver.items", in, b.Operation{})
}
func (c *Client) Settings(ctx context.Context, in b.Empty) (b.ServerConfiguration, error) {
	return b.Call[b.ServerConfiguration](ctx, c.Host, "mediaserver.settings.get", in, b.Operation{})
}
func (c *Client) Refresh(ctx context.Context, in b.ServerRefresh, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "mediaserver.refresh", in, op)
}
func (c *Client) UpdateSettings(ctx context.Context, in b.ServerConfiguration, op b.Operation) (b.ServerConfiguration, error) {
	return b.Call[b.ServerConfiguration](ctx, c.Host, "mediaserver.settings.update", in, op)
}
func (c *Client) Sessions(ctx context.Context, in b.PageRequest) (b.Page[b.PlaybackSession], error) {
	return b.Call[b.Page[b.PlaybackSession]](ctx, c.Host, "mediaserver.sessions", in, b.Operation{})
}
