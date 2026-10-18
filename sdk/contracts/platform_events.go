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
	Phase     string `json:"phase"`
	RequestID string `json:"requestId"`
	Kind      string `json:"kind"`
	ObjectID  string `json:"objectId"`
	Title     string `json:"title,omitempty"`
}
type EffectReceipt struct {
	State string `json:"state"`
}

func ValidEffectRequest(v EffectRequest) bool {
	return OpaqueID(v.RequestID) && OpaqueID(v.ObjectID) && slices.Contains([]string{"notification", "asset", "download"}, v.Kind) && slices.Contains([]string{"", "offer", "execute"}, v.Phase) && len([]rune(v.Title)) <= 200 && !strings.ContainsAny(v.Title, "\x00\r\n")
}
func ValidateEffect(v EffectRequest, r EffectReceipt) error {
	if !ValidEffectRequest(v) {
		return NewError(InvalidExtensionResult)
	}
	if v.Phase == "offer" {
		if r.State == "not_started" || r.State == "ready" {
			return nil
		}
		return NewError(InvalidExtensionResult)
	}
	if v.Phase != "execute" && v.Phase != "" {
		return NewError(InvalidExtensionResult)
	}
	if r.State == "unknown" || r.State == "sent" && v.Kind == "notification" || r.State == "requested" && v.Kind != "notification" {
		return nil
	}
	return NewError(InvalidExtensionResult)
}
func TargetedMessage(t EventType) bool { return t == EventUserMessage || t == EventMessageAction }
