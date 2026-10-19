package business

import (
	"encoding/json"
	"time"
)

type Connection struct {
	ID        string     `json:"id"`
	Kind      string     `json:"kind"`
	Name      string     `json:"name"`
	Enabled   bool       `json:"enabled"`
	Revision  int64      `json:"revision"`
	State     string     `json:"state"`
	CheckedAt *time.Time `json:"checkedAt,omitempty"`
}
type ConnectionCheck struct {
	Passed bool   `json:"passed"`
	State  string `json:"state"`
}
type SettingKeys struct {
	Keys []string `json:"keys"`
}
type SettingsView struct {
	Revision int64                      `json:"revision"`
	Values   map[string]json.RawMessage `json:"values"`
}
type SettingPatch struct {
	Revision int64                      `json:"revision"`
	Values   map[string]json.RawMessage `json:"values"`
}
type NamespaceQuery struct {
	Namespace string `json:"namespace"`
	PageRequest
}
type CacheSummary struct {
	Namespace string `json:"namespace"`
	Entries   int    `json:"entries"`
	Truncated bool   `json:"truncated"`
}
type CacheEntry struct {
	ID         string `json:"id"`
	Bytes      int64  `json:"bytes"`
	TTLSeconds int64  `json:"ttlSeconds"`
	State      string `json:"state"`
}
type CachePage struct {
	Page[CacheEntry]
	Truncated bool `json:"truncated"`
}
type CacheInvalidation struct {
	Namespace string   `json:"namespace"`
	IDs       []string `json:"ids"`
}
type Count struct {
	Count int `json:"count"`
}
type JobQuery struct {
	PageRequest
	Kind  string `json:"kind,omitempty"`
	State string `json:"state,omitempty"`
}
type OwnedJob struct {
	ID          string    `json:"id"`
	Kind        string    `json:"kind"`
	State       string    `json:"state"`
	Attempts    int       `json:"attempts"`
	CreatedAt   time.Time `json:"createdAt"`
	AvailableAt time.Time `json:"availableAt"`
	Cancelable  bool      `json:"cancelable"`
}
type SystemInfo struct {
	Version         string `json:"version"`
	Revision        string `json:"revision"`
	ContractVersion int    `json:"contractVersion"`
}
type HealthSummary struct {
	Postgres  string    `json:"postgres"`
	Redis     string    `json:"redis"`
	Migration *int64    `json:"migration"`
	SampledAt time.Time `json:"sampledAt"`
}
type SystemModule struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Stage string `json:"stage"`
}
type AlertSummary struct {
	Code       string     `json:"code"`
	State      string     `json:"state"`
	Hits       int        `json:"hits"`
	Value      *float64   `json:"value"`
	LastSeenAt *time.Time `json:"lastSeenAt,omitempty"`
}
type Notice struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type DeliverySummary struct {
	ID            string     `json:"id"`
	EventID       string     `json:"eventId"`
	State         string     `json:"state"`
	Attempts      int        `json:"attempts"`
	CreatedAt     time.Time  `json:"createdAt"`
	NextAttemptAt *time.Time `json:"nextAttemptAt,omitempty"`
}
type DeliveryRetry struct {
	ID                  string `json:"id"`
	AcceptDuplicateRisk bool   `json:"acceptDuplicateRisk"`
}
type ProxySummary struct {
	Configured bool  `json:"configured"`
	TMDB       bool  `json:"tmdb"`
	Sites      bool  `json:"sites"`
	Revision   int64 `json:"revision"`
}

// HostRequest is public exact-origin GET/HEAD only, with host proxy profile.
// There is no arbitrary header/cookie/auth/body interface or raw response body.
type HostRequest struct {
	Profile string `json:"profile"`
	Origin  string `json:"origin"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}
type HTTPResult struct {
	Status      int    `json:"status"`
	ContentType string `json:"contentType"`
	BytesRead   int64  `json:"bytesRead"`
}

// NotificationPending proves the host notification outbox has accepted the notice; it is not delivery evidence.
const NotificationPending = "pending"
