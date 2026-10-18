package contracts

import (
	"slices"
	"strings"
)

const EventPluginConfigChanged EventType = "plugin.config.changed"

type NotificationQueued struct {
	NotificationID string `json:"notificationId"`
	Origin         string `json:"origin"`
}
type OperationFailed struct {
	Source string `json:"source"`
	Code   string `json:"code"`
	Count  int64  `json:"count"`
	Window string `json:"window"`
}

// Diff keys are an explicit non-secret allowlist; values and configurations are
// never included (even a server name/template may contain administrator secrets).
type SettingsChanged struct {
	Area     string   `json:"area"`
	Revision int64    `json:"revision"`
	Changed  []string `json:"changed"`
}
type PluginConfigChanged struct {
	InstanceID string `json:"instanceId"`
	Revision   int64  `json:"revision"`
}
type InboundMessage struct {
	InstanceID string `json:"instanceId"`
	Command    string `json:"command"`
	Actor      string `json:"actor"`
	Kind       string `json:"kind"`
}
type CommandInput struct {
	Channel   string `json:"channel"`
	Actor     string `json:"actor"`
	Arguments string `json:"arguments"`
}

// External side effects use explicit outcomes, never bool success. not_started
// is a provider's pre-effect decline; any error/timeout after dispatch is unknown.
type EffectRequest struct {
	RequestID string `json:"requestId"`
	Kind      string `json:"kind"`
	ObjectID  string `json:"objectId"`
	Title     string `json:"title,omitempty"`
}
type EffectReceipt struct {
	State string `json:"state"`
}

func ValidEffectRequest(v EffectRequest) bool {
	return OpaqueID(v.RequestID) && OpaqueID(v.ObjectID) && slices.Contains([]string{"notification", "asset", "download"}, v.Kind) && len([]rune(v.Title)) <= 200 && !strings.ContainsAny(v.Title, "\x00\r\n")
}
func ValidateEffect(v EffectRequest, r EffectReceipt) error {
	if !ValidEffectRequest(v) || !slices.Contains([]string{"not_started", "sent", "requested", "unknown"}, r.State) || r.State == "sent" && v.Kind != "notification" || r.State == "requested" && v.Kind == "notification" {
		return NewError(InvalidExtensionResult)
	}
	return nil
}
func TargetedMessage(t EventType) bool { return t == EventUserMessage || t == EventMessageAction }
