// Package mediaprobe exposes safe typed media Host APIs; no playback URL or member credential methods.
package mediaprobe

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Status(ctx context.Context, in b.Empty) (b.ProcessingStatus, error) {
	return b.Call[b.ProcessingStatus](ctx, c.Host, "mediaprobe.status", in, b.Operation{})
}
func (c *Client) List(ctx context.Context, in b.ProcessingQuery) (b.Page[b.ProcessingTask], error) {
	return b.Call[b.Page[b.ProcessingTask]](ctx, c.Host, "mediaprobe.list", in, b.Operation{})
}
func (c *Client) Retry(ctx context.Context, in b.IDs, op b.Operation) (b.RetryCounts, error) {
	return b.Call[b.RetryCounts](ctx, c.Host, "mediaprobe.retry", in, op)
}
func (c *Client) RetryStills(ctx context.Context, in b.IDs, op b.Operation) (b.RetryCounts, error) {
	return b.Call[b.RetryCounts](ctx, c.Host, "mediaprobe.stills.retry", in, op)
}
func (c *Client) GC(ctx context.Context, in b.Empty, op b.Operation) (b.GCCounts, error) {
	return b.Call[b.GCCounts](ctx, c.Host, "mediaprobe.stills.gc", in, op)
}
