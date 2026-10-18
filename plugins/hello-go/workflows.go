package hello

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"time"
)

func (p *Plugin) workflowEvents() error {
	if e := p.events.On(contracts.EventWorkflowExecute, p.onAcquisition); e != nil {
		return e
	}
	for _, typ := range []contracts.ChainEventType{contracts.ChainWorkflowExecution, contracts.ChainAuthIntercept} {
		if e := p.events.OnChainPriority(typ, 10, func(_ context.Context, q events.ChainEnvelope) (json.RawMessage, error) {
			return json.Marshal(contracts.RestrictResult{Allow: true})
		}); e != nil {
			return e
		}
	}
	return nil
}

// Fake provider has no external connection and never issues a real session.
func (p *Plugin) GetAuthProviders(context.Context) ([]sdk.AuthProvider, error) {
	return []sdk.AuthProvider{{ID: "fake", Name: "验收假身份（非登录入口）", Issuer: "hello-test", Audience: "mediatool-test", Origin: "https://hello.invalid", Handler: func(_ context.Context, q contracts.AuthRequest) (contracts.AuthAssertion, error) {
		return contracts.AuthAssertion{Subject: q.Challenge, Issuer: "hello-test", Audience: "mediatool-test", ExpiresAt: time.Now().Add(time.Minute).Unix(), IsAdmin: true}, nil
	}}}, nil
}
func (p *Plugin) f4Check(ctx context.Context, q sdk.ActionContext) (sdk.ActionContext, bool, error) {
	_, e := p.host.SaveData(ctx, "f4-check", map[string]bool{"passed": true, "noOpAction": true, "authSessionIssued": false}, 0)
	if e != nil {
		var old map[string]bool
		rev, read := p.host.GetData(ctx, "f4-check", &old)
		if read != nil {
			return q, false, read
		}
		_, e = p.host.SaveData(ctx, "f4-check", map[string]bool{"passed": true, "noOpAction": true, "authSessionIssued": false}, rev)
	}
	return q, true, e
}
