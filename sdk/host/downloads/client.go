// Package downloads is a typed, credential-free Host client. Host authorization is always live.
package downloads

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(c *host.Client) *Client { return &Client{Host: c} }
func (c *Client) List(ctx context.Context, in b.PageRequest) (b.Page[b.Task], error) {
	return b.Call[b.Page[b.Task]](ctx, c.Host, "downloads.list", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.ID) (b.Task, error) {
	return b.Call[b.Task](ctx, c.Host, "downloads.get", in, b.Operation{})
}
func (c *Client) Progress(ctx context.Context, in b.ID) (b.Task, error) {
	return b.Call[b.Task](ctx, c.Host, "downloads.progress", in, b.Operation{})
}
func (c *Client) Files(ctx context.Context, in b.Selection) (b.Page[b.File], error) {
	return b.Call[b.Page[b.File]](ctx, c.Host, "downloads.files", in, b.Operation{})
}
func (c *Client) Trackers(ctx context.Context, in b.ID) (b.Page[b.Tracker], error) {
	return b.Call[b.Page[b.Tracker]](ctx, c.Host, "downloads.trackers", in, b.Operation{})
}
func (c *Client) ManualList(ctx context.Context, in b.PageRequest) (b.Page[b.Task], error) {
	return b.Call[b.Page[b.Task]](ctx, c.Host, "downloads.manualList", in, b.Operation{})
}
func (c *Client) AddResource(ctx context.Context, in b.AddResource, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.addResource", in, op)
}
func (c *Client) ManualCreate(ctx context.Context, in b.ManualSource, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.manualCreate", in, op)
}
func (c *Client) ManualRetry(ctx context.Context, in b.Retry, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.manualRetry", in, op)
}
func (c *Client) Pause(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.pause", in, op)
}
func (c *Client) Resume(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.resume", in, op)
}
func (c *Client) Retry(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.retry", in, op)
}
func (c *Client) Recheck(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.recheck", in, op)
}
func (c *Client) PreviewCleanup(ctx context.Context, in b.ID) (b.Cleanup, error) {
	return b.Call[b.Cleanup](ctx, c.Host, "downloads.cleanup.preview", in, b.Operation{})
}
func (c *Client) ApproveCleanup(ctx context.Context, in b.DeleteRequest, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.cleanup.approve", in, op)
}
func (c *Client) AllHashes(ctx context.Context, in b.PageRequest) (b.Page[string], error) {
	return b.Call[b.Page[string]](ctx, c.Host, "downloads.downloader.allHashes", in, b.Operation{})
}
func (c *Client) Inspect(ctx context.Context, in b.ID) (b.Task, error) {
	return b.Call[b.Task](ctx, c.Host, "downloads.downloader.inspect", in, b.Operation{})
}
func (c *Client) SetPaused(ctx context.Context, in b.Paused, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.downloader.setPaused", in, op)
}
func (c *Client) DeleteOwned(ctx context.Context, in b.DeleteRequest, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.downloader.deleteOwned", in, op)
}
func (c *Client) ParseTorrent(ctx context.Context, in b.TorrentInput) (b.Torrent, error) {
	return b.Call[b.Torrent](ctx, c.Host, "metadata.torrent.parse", in, b.Operation{})
}

func (c *Client) UploadControl(ctx context.Context, in b.UploadControl, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "downloads.uploadControl", in, op)
}
