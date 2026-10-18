// Package sites is a typed, credential-free Host client. Host authorization is always live.
package sites

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(c *host.Client) *Client { return &Client{Host: c} }
func (c *Client) List(ctx context.Context, in b.PageRequest) (b.Page[b.Site], error) {
	return b.Call[b.Page[b.Site]](ctx, c.Host, "sites.list", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.ID) (b.Site, error) {
	return b.Call[b.Site](ctx, c.Host, "sites.get", in, b.Operation{})
}
func (c *Client) Metrics(ctx context.Context, in b.ID) (b.Metrics, error) {
	return b.Call[b.Metrics](ctx, c.Host, "sites.metrics", in, b.Operation{})
}
func (c *Client) Request(ctx context.Context, in b.SiteRequest, op b.Operation) (b.SiteResponse, error) {
	return b.Call[b.SiteResponse](ctx, c.Host, "sites.request", in, op)
}
func (c *Client) Checkin(ctx context.Context, in b.SiteRequest, op b.Operation) (b.SiteResponse, error) {
	return b.Call[b.SiteResponse](ctx, c.Host, "sites.checkin", in, op)
}
func (c *Client) DiscoverRSS(ctx context.Context, in b.ID) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "sites.rss.discover", in, b.Operation{})
}
func (c *Client) RefreshRSS(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "sites.rss.refresh", in, op)
}
func (c *Client) CollectMetrics(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "sites.metrics.collect", in, op)
}
func (c *Client) Templates(ctx context.Context, in b.PageRequest) (b.Page[b.Template], error) {
	return b.Call[b.Page[b.Template]](ctx, c.Host, "sites.templates.catalog", in, b.Operation{})
}
func (c *Client) LoadedTemplates(ctx context.Context, in b.PageRequest) (b.Page[b.Template], error) {
	return b.Call[b.Page[b.Template]](ctx, c.Host, "sites.templates.loaded", in, b.Operation{})
}
func (c *Client) LoadTemplates(ctx context.Context, in b.IDs, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "sites.templates.load", in, op)
}
