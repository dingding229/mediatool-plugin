package hello

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"time"
)

func (p *Plugin) organizingEvents() error {
	for _, typ := range []contracts.EventType{contracts.EventTransferFailed, contracts.EventSubtitleTransferComplete, contracts.EventSubtitleTransferFailed, contracts.EventHistoryDeleted, contracts.EventDownloadFileDeleted, contracts.EventDownloadDeleted} {
		if e := p.events.On(typ, p.onAcquisition); e != nil {
			return e
		}
	}
	for _, typ := range []contracts.ChainEventType{contracts.ChainTransferRename, contracts.ChainTransferRenameBuild, contracts.ChainTransferIntercept, contracts.ChainTransferOverwriteCheck, contracts.ChainStorageOperSelection} {
		if e := p.events.OnChainPriority(typ, 20, p.organizingPolicy); e != nil {
			return e
		}
	}
	return nil
}
func (p *Plugin) organizingPolicy(ctx context.Context, e events.ChainEnvelope) (json.RawMessage, error) {
	if string(p.config["organizeDelay"]) == "true" {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	if e.Type == contracts.ChainStorageOperSelection {
		var q contracts.StorageSelection
		if contracts.Decode(e.Payload, &q) != nil {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		return json.Marshal(contracts.StorageChoice{Handle: q.Selected})
	}
	var q contracts.TransferPlan
	if contracts.Decode(e.Payload, &q) != nil || contracts.ValidateTransfer(q) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if e.Type == contracts.ChainTransferRename || e.Type == contracts.ChainTransferRenameBuild {
		var folder string
		if err := json.Unmarshal(p.config["organizeFolder"], &folder); err == nil && folder != "" {
			for i := range q.Paths {
				q.Paths[i] = folder + "/" + q.Paths[i]
			}
		}
		return json.Marshal(contracts.TransferPaths{Paths: q.Paths})
	}
	decision := "pass"
	if string(p.config["organizeReject"]) == "true" {
		decision = "reject"
	}
	return json.Marshal(contracts.TransferDecision{Decision: decision})
}
