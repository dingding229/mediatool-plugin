// Package subscriptions is a typed, credential-free Host client. Host authorization is always live.
package subscriptions

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(c *host.Client) *Client { return &Client{Host: c} }
func (c *Client) List(ctx context.Context, in b.PageRequest) (b.Page[b.Subscription], error) {
	return b.Call[b.Page[b.Subscription]](ctx, c.Host, "subscriptions.list", in, b.Operation{})
}
func (c *Client) Get(ctx context.Context, in b.ID) (b.Subscription, error) {
	return b.Call[b.Subscription](ctx, c.Host, "subscriptions.get", in, b.Operation{})
}
func (c *Client) History(ctx context.Context, in b.Selection) (b.Page[b.History], error) {
	return b.Call[b.Page[b.History]](ctx, c.Host, "subscriptions.history", in, b.Operation{})
}
func (c *Client) Create(ctx context.Context, in b.SubscriptionInput, op b.Operation) (b.Subscription, error) {
	return b.Call[b.Subscription](ctx, c.Host, "subscriptions.create", in, op)
}
func (c *Client) Update(ctx context.Context, in b.SubscriptionInput, op b.Operation) (b.Subscription, error) {
	return b.Call[b.Subscription](ctx, c.Host, "subscriptions.update", in, op)
}
func (c *Client) Queue(ctx context.Context, in b.Queue, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "subscriptions.queue", in, op)
}
func (c *Client) Delete(ctx context.Context, in b.ID, op b.Operation) (b.Receipt, error) {
	return b.Call[b.Receipt](ctx, c.Host, "subscriptions.delete", in, op)
}
func (c *Client) EvaluateRules(ctx context.Context, in b.Evaluation) (b.Evaluated, error) {
	return b.Call[b.Evaluated](ctx, c.Host, "subscriptions.evaluateRules", in, b.Operation{})
}
func (c *Client) Coverage(ctx context.Context, in b.CoverageInput) (b.Coverage, error) {
	return b.Call[b.Coverage](ctx, c.Host, "subscriptions.coverage", in, b.Operation{})
}
func (c *Client) Completion(ctx context.Context, in b.ID) (contracts.CoverageObservation, error) {
	return b.Call[contracts.CoverageObservation](ctx, c.Host, "subscriptions.completion", in, b.Operation{})
}
