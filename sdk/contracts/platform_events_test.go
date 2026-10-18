package contracts_test

import (
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestPlatformEventsAndEffectGolden(t *testing.T) {
	q := contracts.EffectRequest{Phase: "offer", RequestID: aid, ObjectID: bid, Kind: "notification"}
	for _, state := range []string{"not_started", "ready"} {
		if contracts.ValidateEffect(q, contracts.EffectReceipt{State: state}) != nil {
			t.Fatal(state)
		}
	}
	for _, state := range []string{"sent", "unknown", "requested", "delivered_maybe", "success"} {
		if contracts.ValidateEffect(q, contracts.EffectReceipt{State: state}) == nil {
			t.Fatal("ambiguous effect", state)
		}
	}
	execute := q
	execute.Phase = "execute"
	if contracts.ValidateEffect(execute, contracts.EffectReceipt{State: "not_started"}) == nil {
		t.Fatal("post-dispatch decline must never permit another sender")
	}
	if contracts.ValidateEffect(execute, contracts.EffectReceipt{State: "sent"}) != nil {
		t.Fatal("send receipt")
	}
	if !contracts.KnownEventType(contracts.EventPluginConfigChanged) || !contracts.KnownScope("events.subscribe.plugin.config.changed") {
		t.Fatal("missing separate config contract")
	}
	plugintest.Golden(t, "testdata/platform-events.json", map[string]any{
		"notification.queued":       contracts.NotificationQueued{NotificationID: aid, Origin: "builtin"},
		"system.operation.failed":   contracts.OperationFailed{Source: "jobs", Code: "terminal_failure", Count: 3, Window: "2026-10-18T00:00:00Z"},
		"settings.changed":          contracts.SettingsChanged{Area: "native", Revision: 2, Changed: []string{"enabled"}},
		"plugin.config.changed":     contracts.PluginConfigChanged{InstanceID: aid, Revision: 2},
		"messaging.inbound":         contracts.InboundMessage{InstanceID: aid, Command: "hello", Actor: "opaque-actor", Kind: "command"},
		"messaging.action.received": contracts.InboundMessage{InstanceID: aid, Command: "hello", Actor: "opaque-actor", Kind: "action"},
		"command.input":             contracts.CommandInput{Channel: "telegram", Actor: "opaque-actor", Arguments: "demo"},
		"effect.request":            q, "effect.receipt": contracts.EffectReceipt{State: "not_started"},
	})
}
