package workflows

import (
	"context"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
)

type Client struct{ Host *host.Client }

func New(h *host.Client) *Client { return &Client{h} }
func (c *Client) Run(ctx context.Context, q contracts.WorkflowRequest, op b.Operation) (contracts.WorkflowRun, error) {
	return b.Call[contracts.WorkflowRun](ctx, c.Host, "workflows.run", q, op)
}
func (c *Client) Status(ctx context.Context, q contracts.WorkflowRequest, op b.Operation) (contracts.WorkflowRun, error) {
	return b.Call[contracts.WorkflowRun](ctx, c.Host, "workflows.status", q, op)
}
func (c *Client) Cancel(ctx context.Context, q contracts.WorkflowRequest, op b.Operation) (contracts.WorkflowRun, error) {
	return b.Call[contracts.WorkflowRun](ctx, c.Host, "workflows.cancel", q, op)
}
