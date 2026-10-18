package contracts_test

import (
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestPlatformEventsAndEffectGolden(t *testing.T) {
	q := contracts.EffectRequest{RequestID: aid, ObjectID: bid, Kind: "notification"}
	for _, state := range []string{"not_started", "sent", "unknown"} {
		if contracts.ValidateEffect(q, contracts.EffectReceipt{State: state}) != nil {
			t.Fatal(state)
		}
	}
	for _, state := range []string{"requested", "delivered_maybe", "success"} {
		if contracts.ValidateEffect(q, contracts.EffectReceipt{State: state}) == nil {
			t.Fatal("ambiguous effect", state)
		}
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
