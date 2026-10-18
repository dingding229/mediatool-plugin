// Package storage exposes typed, credential-free media/file Host APIs.
package storage

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Roots(ctx context.Context, in b.PageRequest) (b.Page[b.Root], error) {
	return b.Call[b.Page[b.Root]](ctx, c.Host, "storage.roots", in, b.Operation{})
}
func (c *Client) List(ctx context.Context, in b.FileQuery) (b.FilePage, error) {
	return b.Call[b.FilePage](ctx, c.Host, "storage.list", in, b.Operation{})
}
func (c *Client) Stat(ctx context.Context, in b.FileQuery) (b.MediaFile, error) {
	return b.Call[b.MediaFile](ctx, c.Host, "storage.stat", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.ID) (b.FilePlan, error) {
	return b.Call[b.FilePlan](ctx, c.Host, "storage.get", in, b.Operation{})
}
func (c *Client) Preview(ctx context.Context, in b.FileInput, op b.Operation) (b.FilePlan, error) {
	return b.Call[b.FilePlan](ctx, c.Host, "storage.preview", in, op)
}
func (c *Client) Apply(ctx context.Context, in b.Approval, op b.Operation) (b.FilePlan, error) {
	return b.Call[b.FilePlan](ctx, c.Host, "storage.apply", in, op)
}
func (c *Client) PreviewBatch(ctx context.Context, in b.FileInputs, op b.Operation) (b.FileBatch, error) {
	return b.Call[b.FileBatch](ctx, c.Host, "storage.previewBatch", in, op)
}
func (c *Client) ApproveBatch(ctx context.Context, in b.Approval, op b.Operation) (b.FileBatch, error) {
	return b.Call[b.FileBatch](ctx, c.Host, "storage.approveBatch", in, op)
}
func (c *Client) Mkdir(ctx context.Context, in b.FileInput, op b.Operation) (b.FilePlan, error) {
	return b.Call[b.FilePlan](ctx, c.Host, "storage.mkdir", in, op)
}
func (c *Client) Relocate(ctx context.Context, in b.FileInput, op b.Operation) (b.FilePlan, error) {
	return b.Call[b.FilePlan](ctx, c.Host, "storage.relocate", in, op)
}
