package hello

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
)

func (p *Plugin) GetAgentTools(context.Context) ([]sdk.Tool, error) {
	schema := json.RawMessage(`{"type":"object","required":["message"],"properties":{"message":{"type":"string","maxLength":200}},"additionalProperties":false}`)
	return []sdk.Tool{
		{Name: "hello.echo", Description: "原样回显有界文本；无业务副作用", InputSchema: schema, EffectKind: "read", Handler: func(_ context.Context, b json.RawMessage) (json.RawMessage, error) { return b, nil }},
		{Name: "hello.remember", Description: "只更新本实例 tool-message KV，不操作媒体或外部业务", InputSchema: schema, EffectKind: "write", RequiredScopes: []string{"plugin.data.read", "plugin.data.write"}, Handler: func(ctx context.Context, b json.RawMessage) (json.RawMessage, error) {
			var old any
			revision, e := p.host.GetData(ctx, "tool-message", &old)
			var rpc *contracts.Error
			if e != nil && (!errors.As(e, &rpc) || rpc.Code != contracts.NotFound) {
				return nil, e
			}
			if _, e = p.host.SaveData(ctx, "tool-message", b, revision); e != nil {
				return nil, e
			}
			return json.RawMessage(`{"saved":true}`), nil
		}},
	}, nil
}
