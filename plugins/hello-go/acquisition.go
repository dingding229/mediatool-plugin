package hello

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"time"
)

func (p *Plugin) acquisitionEvents() error {
	for _, typ := range []contracts.EventType{contracts.EventSiteUpdated, contracts.EventSiteDeleted, contracts.EventSiteRefreshed, contracts.EventSubscribeAdded, contracts.EventSubscribeModified, contracts.EventSubscribeDeleted} {
		if e := p.events.On(typ, p.onAcquisition); e != nil {
			return e
		}
	}
	if e := p.events.OnChainPriority(contracts.ChainResourceSelection, 20, p.selectResources); e != nil {
		return e
	}
	return p.events.OnChainPriority(contracts.ChainResourceDownload, 20, p.downloadResource)
}
func (p *Plugin) onAcquisition(ctx context.Context, e events.Envelope) error {
	// The UUID key makes replay idempotent. Save only this already-safe typed fact.
	var old json.RawMessage
	if _, err := p.host.GetData(ctx, "acquisition-"+e.ID, &old); err == nil {
		return nil
	}
	_, err := p.host.SaveData(ctx, "acquisition-"+e.ID, e, 0)
	return err
}
func (p *Plugin) resourceDelay(ctx context.Context) error {
	if string(p.config["resourceDelay"]) != "true" {
		return nil
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Second):
		return nil
	}
}
func (p *Plugin) selectResources(ctx context.Context, e events.ChainEnvelope) (json.RawMessage, error) {
	var v contracts.ResourceSelectionRequest
	if contracts.Decode(e.Payload, &v) != nil || contracts.ValidateSelectionRequest(v) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if err := p.resourceDelay(ctx); err != nil {
		return nil, err
	}
	result := contracts.ResourceSelectionResult{IDs: []string{}}
	for _, c := range v.Candidates {
		result.IDs = append(result.IDs, c.ID)
	}
	// Default is a no-op. Optional reverse demonstrates sorting only known IDs.
	if string(p.config["resourceReverse"]) == "true" {
		for i, j := 0, len(result.IDs)-1; i < j; i, j = i+1, j-1 {
			result.IDs[i], result.IDs[j] = result.IDs[j], result.IDs[i]
		}
	}
	return json.Marshal(result)
}
func (p *Plugin) downloadResource(ctx context.Context, e events.ChainEnvelope) (json.RawMessage, error) {
	var v contracts.ResourceDownloadRequest
	if contracts.Decode(e.Payload, &v) != nil || contracts.ValidateDownloadRequest(v) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if err := p.resourceDelay(ctx); err != nil {
		return nil, err
	}
	return json.Marshal(contracts.ResourceDownloadResult{Allow: string(p.config["resourceReject"]) != "true", CleanupAfterImport: v.CleanupAfterImport && string(p.config["resourceNoCleanup"]) != "true"})
}
