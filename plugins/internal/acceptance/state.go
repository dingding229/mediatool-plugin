// Package acceptance shares bounded, credential-free plugin state helpers.
package acceptance

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/notifications"
)

func Load[T any](ctx context.Context, h *host.Client, key string) (T, int64, error) {
	var v T
	rev, err := h.GetData(ctx, key, &v)
	var rpc *contracts.Error
	if errors.As(err, &rpc) && rpc.Code == contracts.NotFound {
		return v, 0, nil
	}
	return v, rev, err
}
func RequireAction(ctx context.Context, h *host.Client, key, code string) error {
	// Only fixed reason codes, never upstream errors, credentials or response text.
	_, rev, err := Load[map[string]string](ctx, h, "status")
	if err != nil {
		return err
	}
	if _, err = h.SaveData(ctx, "status", map[string]string{"state": "needs_action", "code": code}, rev); err != nil {
		return err
	}
	_, notifyErr := notifications.New(h).Send(ctx, b.Notice{Title: "插件需要管理员处理", Text: code + "；操作不会自动重放，请检查实例状态。"}, b.Operation{Key: key, Confirm: true})
	// Preserve needs_action even when notification cannot be queued. No claim of delivery.
	return errors.Join(contracts.NewError(contracts.NeedsAction), notifyErr)
}
func Prop(v any) json.RawMessage { b, _ := json.Marshal(v); return b }
func Field(component, model, label string) contracts.Component {
	return contracts.Component{Component: component, Model: model, Props: map[string]json.RawMessage{"label": Prop(label)}}
}
