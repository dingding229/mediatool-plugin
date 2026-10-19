// Package plugintest provides an in-memory, deny-by-default Host for SDK authors.
// No database, Redis, credentials, network upstreams or filesystem mutations occur.
package plugintest

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host/business"
	"sync"
	"time"
)

type Host struct {
	mu       sync.Mutex
	scopes   map[string]bool
	values   map[string]contracts.Value
	config   contracts.Config
	path     string
	calls    []string
	messages []contracts.Message
	failures map[string]failure
	receipts map[string]businessReceipt
	business map[string]func(context.Context, json.RawMessage) (any, *contracts.Error)
}
type businessReceipt struct {
	request string
	result  json.RawMessage
	code    contracts.Code
}
type failure struct {
	code  contracts.Code
	delay time.Duration
}

func NewHost(scopes []string) *Host {
	h := &Host{scopes: map[string]bool{}, values: map[string]contracts.Value{}, config: contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}, path: "/fixture/plugins/instance/data", failures: map[string]failure{}}
	for _, s := range scopes {
		h.scopes[s] = true
	}
	return h
}
func (h *Host) Revoke(scope string) { h.mu.Lock(); defer h.mu.Unlock(); delete(h.scopes, scope) }
func (h *Host) InjectOnce(method string, code contracts.Code, delay time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failures[method] = failure{code, delay}
}
func (h *Host) Calls() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]string{}, h.calls...)
}
func (h *Host) Messages() []contracts.Message {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]contracts.Message{}, h.messages...)
}
func methodScope(method string) string {
	switch method {
	case "plugin.data.get":
		return "plugin.data.read"
	case "plugin.data.put", "plugin.data.delete":
		return "plugin.data.write"
	case "plugin.data.path":
		return "fs.data"
	case "plugin.config.get":
		return "plugin.config.read"
	case "plugin.config.update":
		return "plugin.config.write"
	case "notifications.send":
		return "notifications.send"
	}
	return ""
}
func (h *Host) Handle(ctx context.Context, b json.RawMessage) (any, *contracts.Error) {
	var call contracts.HostCall
	if contracts.Decode(b, &call) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	h.mu.Lock()
	if len(h.calls) >= 4096 {
		h.mu.Unlock()
		return nil, contracts.NewError(contracts.Busy)
	}
	h.calls = append(h.calls, call.Method)
	f, hasFailure := h.failures[call.Method]
	delete(h.failures, call.Method)
	h.mu.Unlock()
	if hasFailure {
		timer := time.NewTimer(f.delay)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return nil, contracts.NewError(contracts.Cancelled)
		case <-timer.C:
		}
		if f.code != 0 {
			return nil, contracts.NewError(f.code)
		}
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	scope := methodScope(call.Method)
	spec, isBusiness := business.Lookup(call.Method)
	if isBusiness {
		scope = spec.Scope
	}
	if scope == "" {
		return nil, contracts.NewError(contracts.MethodNotFound)
	}
	if !h.scopes[scope] {
		return nil, contracts.NewError(contracts.PermissionDenied)
	}
	if isBusiness {
		for _, required := range business.RequiredScopes(call.Method) {
			if !h.scopes[required] {
				return nil, contracts.NewError(contracts.PermissionDenied)
			}
		}
		if !spec.Supported {
			return nil, contracts.NewError(contracts.Unsupported)
		}
		fn := h.business[call.Method]
		if fn == nil {
			return nil, contracts.NewError(contracts.Unsupported)
		}
		if spec.Level == "R" {
			return fn(ctx, call.Args)
		}
		if len(call.OperationKey) < 8 || len(call.OperationKey) > 100 {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		if h.receipts == nil {
			h.receipts = map[string]businessReceipt{}
		}
		key := call.Method + ":" + call.OperationKey
		if old, ok := h.receipts[key]; ok {
			if old.request != string(call.Args) {
				return nil, contracts.NewError(contracts.Conflict)
			}
			if old.code != 0 {
				return nil, contracts.NewError(old.code)
			}
			return append(json.RawMessage{}, old.result...), nil
		}
		if len(h.receipts) >= 1024 {
			return nil, contracts.NewError(contracts.Busy)
		}
		h.receipts[key] = businessReceipt{request: string(call.Args), code: contracts.ResultUnknown}
		v, rpc := fn(ctx, call.Args)
		if rpc != nil {
			h.receipts[key] = businessReceipt{request: string(call.Args), code: rpc.Code}
			return nil, rpc
		}
		raw, e := json.Marshal(v)
		if e != nil {
			return nil, contracts.NewError(contracts.InternalError)
		}
		h.receipts[key] = businessReceipt{request: string(call.Args), result: raw}
		return v, nil
	}
	var key struct {
		Key string `json:"key"`
	}
	var put struct {
		Key   string          `json:"key"`
		Value json.RawMessage `json:"value"`
	}
	switch call.Method {
	case "plugin.data.get", "plugin.data.delete":
		if contracts.Decode(call.Args, &key) != nil || key.Key == "" || len(key.Key) > 128 {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		val, ok := h.values[key.Key]
		if !ok {
			return nil, contracts.NewError(contracts.NotFound)
		}
		if call.Method == "plugin.data.get" {
			return val, nil
		}
		if call.ExpectedRevision == nil || *call.ExpectedRevision != val.Revision {
			return nil, contracts.NewError(contracts.Conflict)
		}
		delete(h.values, key.Key)
		return struct{}{}, nil
	case "plugin.data.put":
		if contracts.Decode(call.Args, &put) != nil || put.Key == "" || len(put.Key) > 128 || len(put.Value) > 128<<10 {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		val := h.values[put.Key]
		if call.ExpectedRevision == nil || *call.ExpectedRevision != val.Revision {
			return nil, contracts.NewError(contracts.Conflict)
		}
		val = contracts.Value{Value: append(json.RawMessage{}, put.Value...), Revision: val.Revision + 1}
		h.values[put.Key] = val
		return contracts.Revision{Revision: val.Revision}, nil
	case "plugin.config.get":
		return h.config, nil
	case "plugin.config.update":
		var config struct {
			Value json.RawMessage `json:"value"`
		}
		if contracts.Decode(call.Args, &config) != nil {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		if call.ExpectedRevision == nil || *call.ExpectedRevision != h.config.Revision {
			return nil, contracts.NewError(contracts.Conflict)
		}
		h.config = contracts.Config{Value: append(json.RawMessage{}, config.Value...), Revision: h.config.Revision + 1}
		return contracts.Revision{Revision: h.config.Revision}, nil
	case "plugin.data.path":
		return struct {
			Path string `json:"path"`
		}{h.path}, nil
	case "notifications.send":
		var m contracts.Message
		if contracts.Decode(call.Args, &m) != nil {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		h.messages = append(h.messages, m)
		return struct{}{}, nil
	}
	return nil, contracts.NewError(contracts.MethodNotFound)
}

// OnBusiness registers deterministic fixture behavior, not a fabricated business implementation.
// The mock still checks the shared per-method scope before invoking it.
func (h *Host) OnBusiness(method string, fn func(context.Context, json.RawMessage) (any, *contracts.Error)) error {
	if _, ok := business.Lookup(method); !ok || fn == nil {
		return contracts.NewError(contracts.InvalidParams)
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.business == nil {
		h.business = map[string]func(context.Context, json.RawMessage) (any, *contracts.Error){}
	}
	h.business[method] = fn
	return nil
}
