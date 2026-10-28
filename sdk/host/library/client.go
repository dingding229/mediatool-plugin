// Package library exposes typed, credential-free media/file Host APIs.
package library

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Assets(ctx context.Context, in b.PageRequest) (b.Page[b.Asset], error) {
	return b.Call[b.Page[b.Asset]](ctx, c.Host, "library.assets", in, b.Operation{})
}
func (c *Client) Candidates(ctx context.Context, in b.LibraryPage) (b.Page[b.Candidate], error) {
	return b.Call[b.Page[b.Candidate]](ctx, c.Host, "library.candidates", in, b.Operation{})
}
func (c *Client) History(ctx context.Context, in b.LibraryPage) (b.Page[b.LibraryHistory], error) {
	return b.Call[b.Page[b.LibraryHistory]](ctx, c.Host, "library.history", in, b.Operation{})
}
func (c *Client) Scans(ctx context.Context, in b.PageRequest) (b.Page[b.Scan], error) {
	return b.Call[b.Page[b.Scan]](ctx, c.Host, "library.scans", in, b.Operation{})
}
func (c *Client) Recognize(ctx context.Context, in b.ID) (b.Recognition, error) {
	return b.Call[b.Recognition](ctx, c.Host, "library.recognize", in, b.Operation{})
}
func (c *Client) Suggest(ctx context.Context, in b.ID) ([]b.Media, error) {
	return b.Call[[]b.Media](ctx, c.Host, "library.suggest", in, b.Operation{})
}
func (c *Client) OrganizePreview(ctx context.Context, in b.ID, op b.Operation) (b.Recognition, error) {
	return b.Call[b.Recognition](ctx, c.Host, "library.organize.preview", in, op)
}
func (c *Client) CorrectPreview(ctx context.Context, in b.ConfirmIdentity, op b.Operation) (b.IdentityPlan, error) {
	return b.Call[b.IdentityPlan](ctx, c.Host, "library.correct.preview", in, op)
}
func (c *Client) StartScan(ctx context.Context, in b.ScanStart, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "library.scan.start", in, op)
}
func (c *Client) ControlScan(ctx context.Context, in b.Control, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "library.scan.control", in, op)
}
func (c *Client) Organize(ctx context.Context, in b.ConfirmIdentity, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "library.organize.confirm", in, op)
}
func (c *Client) Correct(ctx context.Context, in b.ConfirmIdentity, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "library.correct.apply", in, op)
}
func (c *Client) Refresh(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "library.refresh", in, op)
}
func (c *Client) HideHistory(ctx context.Context, in b.HistorySelections, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "library.history.hide", in, op)
}
