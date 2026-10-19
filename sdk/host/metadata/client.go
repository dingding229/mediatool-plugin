// Package metadata exposes typed, credential-free media/file Host APIs.
package metadata

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Parse(ctx context.Context, in b.Text) (b.NameHint, error) {
	return b.Call[b.NameHint](ctx, c.Host, "metadata.parse", in, b.Operation{})
}
func (c *Client) Attributes(ctx context.Context, in b.Text) (map[string]string, error) {
	return b.Call[map[string]string](ctx, c.Host, "metadata.attributes.parse", in, b.Operation{})
}
func (c *Client) ValidateAttributes(ctx context.Context, in b.Attribute) (bool, error) {
	return b.Call[bool](ctx, c.Host, "metadata.attributes.validate", in, b.Operation{})
}
func (c *Client) Render(ctx context.Context, in b.RenderName) (string, error) {
	return b.Call[string](ctx, c.Host, "metadata.name.render", in, b.Operation{})
}
func (c *Client) Episodes(ctx context.Context, in b.Text) (b.EpisodeCoverage, error) {
	return b.Call[b.EpisodeCoverage](ctx, c.Host, "metadata.episodes", in, b.Operation{})
}
func (c *Client) Normalize(ctx context.Context, in b.EpisodeCoverage) (b.EpisodeCoverage, error) {
	return b.Call[b.EpisodeCoverage](ctx, c.Host, "metadata.normalize", in, b.Operation{})
}
func (c *Client) Search(ctx context.Context, in b.MetadataQuery) ([]b.Media, error) {
	return b.Call[[]b.Media](ctx, c.Host, "metadata.search", in, b.Operation{})
}
func (c *Client) Details(ctx context.Context, in b.MetadataQuery) (b.Media, error) {
	return b.Call[b.Media](ctx, c.Host, "metadata.details", in, b.Operation{})
}
func (c *Client) StoredDetails(ctx context.Context, in b.MetadataQuery) (b.StoredMetadata, error) {
	return b.Call[b.StoredMetadata](ctx, c.Host, "metadata.storedDetails", in, b.Operation{})
}
func (c *Client) Refresh(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "metadata.refresh", in, op)
}
func (c *Client) NativeDetails(ctx context.Context, in b.NativeTarget) (b.NativeMetadata, error) {
	return b.Call[b.NativeMetadata](ctx, c.Host, "metadata.nativeDetails", in, b.Operation{})
}
func (c *Client) RefreshDetails(ctx context.Context, in b.NativeTarget, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "metadata.refreshDetails", in, op)
}

func (c *Client) Wishlist(ctx context.Context, in b.WishlistRequest) (b.WishlistPage, error) {
	return b.Call[b.WishlistPage](ctx, c.Host, "metadata.wishlist", in, b.Operation{})
}
