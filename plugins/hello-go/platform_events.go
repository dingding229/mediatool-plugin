package hello

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
)

func (p *Plugin) platformEvents() error {
	for _, typ := range []contracts.EventType{contracts.EventNoticeMessage, contracts.EventSystemError, contracts.EventConfigChanged, contracts.EventPluginConfigChanged, contracts.EventUserMessage, contracts.EventMessageAction} {
		if e := p.events.On(typ, p.onAcquisition); e != nil {
			return e
		}
	}
	return nil
}
func (p *Plugin) platformModules() []sdk.Module {
	out := []sdk.Module{}
	for _, id := range []string{"notifications.channel", "mediaserver.refresh"} {
		out = append(out, sdk.Module{ID: id, Priority: 10, Handler: func(_ context.Context, b json.RawMessage) (json.RawMessage, error) {
			var q contracts.EffectRequest
			if contracts.Decode(b, &q) != nil || !contracts.ValidEffectRequest(q) {
				return nil, contracts.NewError(contracts.InvalidParams)
			}
			state := "unknown"
			if q.Phase == "offer" {
				state = "not_started"
			}
			return json.Marshal(contracts.EffectReceipt{State: state})
		}})
	}
	return out
}
