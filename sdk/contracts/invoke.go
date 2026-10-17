package contracts

import "encoding/json"

// Descriptor is the immutable wire metadata; Go handlers remain worker-local.
const DescriptorVersion = 2

type Descriptor struct {
	Kind            string          `json:"kind"`
	ID              string          `json:"id"`
	Handler         string          `json:"handler"`
	RequiredScopes  []string        `json:"requiredScopes"`
	Name            string          `json:"name,omitempty"`
	Path            string          `json:"path,omitempty"`
	Methods         []string        `json:"methods,omitempty"`
	Auth            string          `json:"auth,omitempty"`
	Summary         string          `json:"summary,omitempty"`
	Description     string          `json:"description,omitempty"`
	Trigger         string          `json:"trigger,omitempty"`
	Cron            string          `json:"cron,omitempty"`
	IntervalMS      int64           `json:"intervalMs,omitempty"`
	At              string          `json:"at,omitempty"`
	ScheduleSetting string          `json:"scheduleSetting,omitempty"`
	Kwargs          json.RawMessage `json:"kwargs,omitempty"`
	Cmd             string          `json:"cmd,omitempty"`
	Category        string          `json:"category,omitempty"`
	Data            json.RawMessage `json:"data,omitempty"`
	Priority        int             `json:"priority,omitempty"`
}
type Snapshot struct {
	Version     int          `json:"version"`
	Descriptors []Descriptor `json:"descriptors"`
}
type Invocation struct {
	Handler string          `json:"handler"`
	Args    json.RawMessage `json:"args"`
}
type APIRequest struct {
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   map[string][]string `json:"query,omitempty"`
	Headers map[string]string   `json:"headers,omitempty"`
	Body    []byte              `json:"body,omitempty"`
}
type APIResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    []byte            `json:"body,omitempty"`
}
type ActionContext struct {
	Data map[string]json.RawMessage `json:"data"`
}
type ActionResult struct {
	Context  ActionContext `json:"context"`
	Continue bool          `json:"continue"`
}
