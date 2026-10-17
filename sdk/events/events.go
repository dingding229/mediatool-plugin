// Package events implements typed, local registrations; Host event production is F2/F4.
package events

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"sort"
	"sync"
)

type EventType = contracts.EventType
type ChainEventType = contracts.ChainEventType
type Envelope struct {
	ID            string          `json:"id"`
	Type          EventType       `json:"type"`
	SchemaVersion int             `json:"schemaVersion"`
	Payload       json.RawMessage `json:"payload"`
}
type ChainEnvelope struct {
	Type          ChainEventType  `json:"type"`
	SchemaVersion int             `json:"schemaVersion"`
	Payload       json.RawMessage `json:"payload"`
}
type Handler func(context.Context, Envelope) error
type ChainHandler func(context.Context, ChainEnvelope) (json.RawMessage, error)
type Descriptor struct {
	Type     string `json:"type"`
	Kind     string `json:"kind"`
	Priority int    `json:"priority,omitempty"`
}
type Registry struct {
	mu         sync.RWMutex
	handlers   map[EventType]Handler
	chains     map[ChainEventType]ChainHandler
	priorities map[ChainEventType]int
	frozen     bool
}

func New() *Registry {
	return &Registry{handlers: map[EventType]Handler{}, chains: map[ChainEventType]ChainHandler{}, priorities: map[ChainEventType]int{}}
}
func (r *Registry) On(t EventType, h Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.frozen || !contracts.KnownEventType(t) || h == nil || r.handlers[t] != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	if len(r.handlers)+len(r.chains) >= 64 {
		return contracts.NewError(contracts.Busy)
	}
	r.handlers[t] = h
	return nil
}
func (r *Registry) OnChain(t ChainEventType, h ChainHandler) error { return r.OnChainPriority(t, 0, h) }
func (r *Registry) OnChainPriority(t ChainEventType, priority int, h ChainHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if priority < -1000 || priority > 1000 || r.frozen || !contracts.KnownChainEventType(t) || h == nil || r.chains[t] != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	if len(r.handlers)+len(r.chains) >= 64 {
		return contracts.NewError(contracts.Busy)
	}
	r.chains[t] = h
	r.priorities[t] = priority
	return nil
}
func (r *Registry) Freeze() { r.mu.Lock(); defer r.mu.Unlock(); r.frozen = true }
func (r *Registry) Descriptors() []Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Descriptor{}
	for t := range r.handlers {
		out = append(out, Descriptor{Type: string(t), Kind: "event"})
	}
	for t := range r.chains {
		out = append(out, Descriptor{Type: string(t), Kind: "chain", Priority: r.priorities[t]})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind == out[j].Kind {
			return out[i].Type < out[j].Type
		}
		return out[i].Kind < out[j].Kind
	})
	return out
}
func (r *Registry) Dispatch(ctx context.Context, e Envelope) error {
	r.mu.RLock()
	h := r.handlers[e.Type]
	r.mu.RUnlock()
	if h == nil || e.SchemaVersion != 1 || contracts.ValidateJSON(e.Payload) != nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	return h(ctx, e)
}
func (r *Registry) DispatchChain(ctx context.Context, e ChainEnvelope) (json.RawMessage, error) {
	r.mu.RLock()
	h := r.chains[e.Type]
	r.mu.RUnlock()
	if h == nil || e.SchemaVersion != 1 || contracts.ValidateJSON(e.Payload) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	out, err := h(ctx, e)
	if err != nil {
		return nil, err
	}
	if contracts.ValidateJSON(out) != nil {
		return nil, contracts.NewError(contracts.InvalidExtensionResult)
	}
	return out, nil
}

// Typed strictly decodes T. Event-specific DTO/schema pairs are added alongside real producers.
func Typed[T any](h func(context.Context, T) error) Handler {
	return func(ctx context.Context, e Envelope) error {
		var payload T
		if err := contracts.Decode(e.Payload, &payload); err != nil {
			return err
		}
		return h(ctx, payload)
	}
}
