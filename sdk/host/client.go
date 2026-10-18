// Package host exposes typed helpers for an instance's own data and configuration.
// It never authenticates a caller: the Host must check scopes/revisions on EVERY call.
package host

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"strings"
	"unicode/utf8"
)

type Caller interface {
	Call(context.Context, string, any, any) error
}
type Client struct{ rpc Caller }

func New(rpc Caller) *Client { return &Client{rpc: rpc} }
func (c *Client) call(ctx context.Context, method string, args any, rev *int64, out any) error {
	if c == nil || c.rpc == nil || ctx == nil || (rev != nil && *rev < 0) {
		return contracts.NewError(contracts.InvalidParams)
	}
	raw, e := json.Marshal(args)
	if e != nil || len(raw) > 128<<10 {
		return contracts.NewError(contracts.InvalidParams)
	}
	if contracts.ValidateJSON(raw) != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	return c.rpc.Call(ctx, "host.call", contracts.HostCall{Method: method, Args: raw, ExpectedRevision: rev}, out)
}
func keyOK(key string) bool {
	return key != "" && utf8.ValidString(key) && len(key) <= 128 && !strings.ContainsAny(key, "\x00\r\n")
}
func (c *Client) SaveData(ctx context.Context, key string, value any, expectedRevision int64) (contracts.Revision, error) {
	var out contracts.Revision
	if !keyOK(key) {
		return out, contracts.NewError(contracts.InvalidParams)
	}
	e := c.call(ctx, "plugin.data.put", struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}{key, value}, &expectedRevision, &out)
	return out, e
}
func (c *Client) GetData(ctx context.Context, key string, out any) (int64, error) {
	if !keyOK(key) || out == nil {
		return 0, contracts.NewError(contracts.InvalidParams)
	}
	var val contracts.Value
	if e := c.call(ctx, "plugin.data.get", struct {
		Key string `json:"key"`
	}{key}, nil, &val); e != nil {
		return 0, e
	}
	return val.Revision, contracts.Decode(val.Value, out)
}
func (c *Client) DelData(ctx context.Context, key string, expectedRevision int64) error {
	if !keyOK(key) {
		return contracts.NewError(contracts.InvalidParams)
	}
	return c.call(ctx, "plugin.data.delete", struct {
		Key string `json:"key"`
	}{key}, &expectedRevision, nil)
}
func (c *Client) GetDataPath(ctx context.Context) (string, error) {
	var out struct {
		Path string `json:"path"`
	}
	e := c.call(ctx, "plugin.data.path", struct{}{}, nil, &out)
	if e == nil && (out.Path == "" || len(out.Path) > 4096) {
		e = contracts.NewError(contracts.InvalidExtensionResult)
	}
	return out.Path, e
}
func (c *Client) GetConfig(ctx context.Context, out any) (int64, error) {
	if out == nil {
		return 0, contracts.NewError(contracts.InvalidParams)
	}
	var val contracts.Config
	if e := c.call(ctx, "plugin.config.get", struct{}{}, nil, &val); e != nil {
		return 0, e
	}
	return val.Revision, contracts.Decode(val.Value, out)
}
func (c *Client) UpdateConfig(ctx context.Context, value any, expectedRevision int64) (contracts.Revision, error) {
	var out contracts.Revision
	e := c.call(ctx, "plugin.config.update", struct {
		Value any `json:"value"`
	}{value}, &expectedRevision, &out)
	return out, e
}
func (c *Client) PostMessage(ctx context.Context, m contracts.Message) error {
	if len(m.Title) > 256 || len(m.Text) > 8192 || len(m.ChannelID) > 128 {
		return contracts.NewError(contracts.InvalidParams)
	}
	return c.call(ctx, "notifications.send", m, nil, nil)
}

// Invoke is used by typed module clients; it grants no authority and exposes no host objects.
func (c *Client) Invoke(ctx context.Context, method string, args any, key string, rev *int64, out any) error {
	if c == nil || c.rpc == nil || ctx == nil || rev != nil && *rev < 0 {
		return contracts.NewError(contracts.InvalidParams)
	}
	raw, e := json.Marshal(args)
	if e != nil || len(raw) > 128<<10 {
		return contracts.NewError(contracts.InvalidParams)
	}
	if e = contracts.ValidateJSON(raw); e != nil {
		return e
	}
	return c.rpc.Call(ctx, "host.call", contracts.HostCall{Method: method, Args: raw, OperationKey: key, ExpectedRevision: rev}, out)
}
