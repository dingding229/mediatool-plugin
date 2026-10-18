// Package mediausers exposes safe typed media Host APIs; no playback URL or member credential methods.
package mediausers

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) List(ctx context.Context, in b.PageRequest) (b.Page[b.ServerUser], error) {
	return b.Call[b.Page[b.ServerUser]](ctx, c.Host, "mediausers.list", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.ID) (b.ServerUser, error) {
	return b.Call[b.ServerUser](ctx, c.Host, "mediausers.get", in, b.Operation{})
}
func (c *Client) Devices(ctx context.Context, in b.Selection) (b.Page[b.ServerDevice], error) {
	return b.Call[b.Page[b.ServerDevice]](ctx, c.Host, "mediausers.devices", in, b.Operation{})
}
func (c *Client) History(ctx context.Context, in b.Selection) (b.Page[b.WatchHistory], error) {
	return b.Call[b.Page[b.WatchHistory]](ctx, c.Host, "mediausers.history", in, b.Operation{})
}
func (c *Client) Update(ctx context.Context, in b.UserPatch, op b.Operation) (b.ServerUser, error) {
	return b.Call[b.ServerUser](ctx, c.Host, "mediausers.update", in, op)
}
func (c *Client) RevokeDevice(ctx context.Context, in b.DeviceRevocation, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "mediausers.revokeDevice", in, op)
}
