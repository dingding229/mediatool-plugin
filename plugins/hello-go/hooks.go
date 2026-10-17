package hello

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"time"
)

func (p *Plugin) recognize(ctx context.Context, e events.ChainEnvelope) (json.RawMessage, error) {
	var v contracts.NameRecognition
	if contracts.Decode(e.Payload, &v) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if string(p.config["chainError"]) == "true" {
		return nil, contracts.NewError(contracts.NeedsAction)
	}
	if string(p.config["chainDelay"]) == "true" {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	var title string
	_ = json.Unmarshal(p.config["recognizeTitle"], &title)
	if title != "" {
		v.Hint.Title = title
	}
	if string(p.config["intercept"]) == "true" {
		v.Continue = false
	}
	return json.Marshal(v)
}
func (p *Plugin) GetModule(context.Context) ([]sdk.Module, error) {
	return []sdk.Module{{ID: "metadata.name_recognize", Priority: 10, Handler: func(_ context.Context, b json.RawMessage) (json.RawMessage, error) {
		var v contracts.NameRecognition
		if contracts.Decode(b, &v) != nil {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		var t string
		_ = json.Unmarshal(p.config["moduleTitle"], &t)
		if t != "" {
			v.Hint.Title = t
		}
		return json.Marshal(v)
	}}}, nil
}
func (p *Plugin) onDownload(ctx context.Context, e events.Envelope) error {
	var old any
	rev, err := p.host.GetData(ctx, "download-"+e.ID, &old)
	var rpc *contracts.Error
	if err == nil {
		return nil
	}
	if !errors.As(err, &rpc) || rpc.Code != contracts.NotFound {
		return err
	}
	_, err = p.host.SaveData(ctx, "download-"+e.ID, json.RawMessage(e.Payload), rev)
	return err
}
func (p *Plugin) GetDashboardMeta(context.Context) ([]sdk.DashboardMeta, error) {
	return []sdk.DashboardMeta{{Key: "starts", Title: "启动次数", Refresh: 30, Border: true, Cols: 6}, {Key: "events", Title: "事件与扩展", Refresh: 60, Border: false, Cols: 6}}, nil
}
func (p *Plugin) GetDashboard(ctx context.Context, q sdk.DashboardRequest) (sdk.Dashboard, error) {
	return p.GetPage(ctx, q)
}
