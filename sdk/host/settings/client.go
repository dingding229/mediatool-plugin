// Package settings provides a bounded, credential-free platform Host client.
package settings

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Get(ctx context.Context, in b.SettingKeys) (b.SettingsView, error) {
	return b.Call[b.SettingsView](ctx, c.Host, "settings.get", in, b.Operation{})
}
func (c *Client) Patch(ctx context.Context, in b.SettingPatch, op b.Operation) (b.SettingsView, error) {
	return b.Call[b.SettingsView](ctx, c.Host, "settings.patch", in, op)
}
