// Package playback exposes safe typed media Host APIs; no playback URL or member credential methods.
package playback

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Sessions(ctx context.Context, in b.PageRequest) (b.Page[b.PlaybackSession], error) {
	return b.Call[b.Page[b.PlaybackSession]](ctx, c.Host, "playback.sessions", in, b.Operation{})
}
func (c *Client) Summary(ctx context.Context, in b.Empty) (b.PlaybackCounts, error) {
	return b.Call[b.PlaybackCounts](ctx, c.Host, "playback.summary", in, b.Operation{})
}
func (c *Client) AccessSummary(ctx context.Context, in b.Empty) (b.AccessCounts, error) {
	return b.Call[b.AccessCounts](ctx, c.Host, "playeraccess.summary", in, b.Operation{})
}
