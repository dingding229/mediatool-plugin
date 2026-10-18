package hello

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"time"
)

func (p *Plugin) discoveryEvents() error {
	for _, typ := range []contracts.EventType{contracts.EventMetadataScrape, contracts.EventSubscribeComplete} {
		if e := p.events.On(typ, p.onAcquisition); e != nil {
			return e
		}
	}
	for _, typ := range []contracts.ChainEventType{contracts.ChainMediaRecognizeConvert, contracts.ChainSubscribeEpisodesRefresh, contracts.ChainSubscribeCompletionCheck, contracts.ChainDiscoverSource, contracts.ChainRecommendSource} {
		if e := p.events.OnChainPriority(typ, 20, p.discoveryPolicy); e != nil {
			return e
		}
	}
	return nil
}
func (p *Plugin) discoveryPolicy(ctx context.Context, e events.ChainEnvelope) (json.RawMessage, error) {
	if string(p.config["discoveryDelay"]) == "true" {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}
	}
	switch e.Type {
	case contracts.ChainMediaRecognizeConvert:
		var q contracts.MediaConversion
		if err := contracts.Decode(e.Payload, &q); err != nil {
			return nil, err
		}
		return json.Marshal(contracts.MediaChoice{Selected: q.Selected})
	case contracts.ChainSubscribeEpisodesRefresh:
		var q contracts.EpisodeExpectation
		if err := contracts.Decode(e.Payload, &q); err != nil {
			return nil, err
		}
		return json.Marshal(contracts.EpisodeSupplement{Expected: q.Expected})
	case contracts.ChainSubscribeCompletionCheck:
		var q contracts.CompletionCheck
		if err := contracts.Decode(e.Payload, &q); err != nil {
			return nil, err
		}
		return json.Marshal(contracts.CompletionDecision{Allow: q.Complete})
	default:
		return json.Marshal(contracts.SourceResult{Items: []contracts.SourceCandidate{}})
	}
}
