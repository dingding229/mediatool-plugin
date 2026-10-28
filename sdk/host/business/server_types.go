package business

import "time"

// ServerQuery is administrative catalog access, not a member identity. Deliberately
// no UserID, VisibleLibraries, URL, headers, device key or playback parameters.
type ServerQuery struct {
	PageRequest
	Provider  string `json:"provider"`
	LibraryID string `json:"libraryId,omitempty"`
	ID        string `json:"id,omitempty"`
	MediaType string `json:"mediaType,omitempty"`
}
type ServerProvider struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Catalog bool   `json:"catalog"`
	Refresh bool   `json:"refresh"`
}
type ServerLibrary struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MediaType string `json:"mediaType"`
	ItemCount int64  `json:"itemCount"`
}
type ServerItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Year      string `json:"year"`
	MediaType string `json:"mediaType"`
	Overview  string `json:"overview"`
	Season    *int   `json:"season,omitempty"`
	Episode   *int   `json:"episode,omitempty"`
}
type ServerRefresh struct {
	Provider string `json:"provider"`
	AssetID  string `json:"assetId"`
}
type ServerConfiguration struct {
	ServerName string `json:"serverName"`
	Enabled    bool   `json:"enabled"`
	Revision   int64  `json:"revision"`
}
type ServerGrant struct {
	LibraryID   string `json:"libraryId"`
	CanPlay     bool   `json:"canPlay"`
	CanDownload bool   `json:"canDownload"`
}
type ServerUser struct {
	ID          string        `json:"id"`
	Username    string        `json:"username"`
	DisplayName string        `json:"displayName"`
	Enabled     bool          `json:"enabled"`
	Revision    int64         `json:"revision"`
	MaxStreams  int           `json:"maxStreams"`
	ExpiresAt   *time.Time    `json:"expiresAt"`
	Grants      []ServerGrant `json:"grants"`
}

// Omitted patch fields retain their current values; explicit ClearExpiry clears
// expiration. Passwords, authentication and new account creation are absent.
type UserPatch struct {
	ID          string         `json:"id"`
	Revision    int64          `json:"revision"`
	DisplayName *string        `json:"displayName,omitempty"`
	Enabled     *bool          `json:"enabled,omitempty"`
	MaxStreams  *int           `json:"maxStreams,omitempty"`
	ExpiresAt   *time.Time     `json:"expiresAt,omitempty"`
	ClearExpiry bool           `json:"clearExpiry,omitempty"`
	Grants      *[]ServerGrant `json:"grants,omitempty"`
}
type ServerDevice struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Client     string    `json:"client"`
	Version    string    `json:"version"`
	LastSeenAt time.Time `json:"lastSeenAt"`
}
type DeviceRevocation struct {
	ID       string `json:"id"`
	DeviceID string `json:"deviceId"`
	Revision int64  `json:"revision"`
}
type WatchHistory struct {
	ItemID        string    `json:"itemId"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	PositionTicks int64     `json:"positionTicks"`
	PlayCount     int       `json:"playCount"`
	Played        bool      `json:"played"`
	LastPlayedAt  time.Time `json:"lastPlayedAt"`
}
type PlaybackSession struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	UserName       string    `json:"userName"`
	DeviceID       string    `json:"deviceId"`
	DeviceName     string    `json:"deviceName"`
	ItemID         string    `json:"itemId"`
	ItemName       string    `json:"itemName"`
	State          string    `json:"state"`
	PositionTicks  int64     `json:"positionTicks"`
	LastProgressAt time.Time `json:"lastProgressAt"`
}
type PlaybackCounts struct {
	Playing int64 `json:"playing"`
	Paused  int64 `json:"paused"`
}
type AccessCounts struct {
	Sampled    int64 `json:"sampled"`
	Allowed    int64 `json:"allowed"`
	Denied     int64 `json:"denied"`
	MaxRecords int   `json:"maxRecords"`
}
type ProcessingQuery struct {
	PageRequest
	Kind  string `json:"kind"`
	State string `json:"state,omitempty"`
}
type ProcessingTask struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	State         string     `json:"state"`
	ErrorReason   string     `json:"errorReason"`
	Attempts      int        `json:"attempts"`
	NextAttemptAt *time.Time `json:"nextAttemptAt"`
}
type ProcessingStateCount struct {
	State string `json:"state"`
	Count int64  `json:"count"`
}
type ProcessingStatus struct {
	Enabled       bool                   `json:"enabled"`
	ProbeEnabled  bool                   `json:"probeEnabled"`
	StillsEnabled bool                   `json:"stillsEnabled"`
	Probe         []ProcessingStateCount `json:"probe"`
	Still         []ProcessingStateCount `json:"still"`
	UsedBytes     int64                  `json:"usedBytes"`
	ReservedBytes int64                  `json:"reservedBytes"`
	LimitBytes    int64                  `json:"limitBytes"`
	UnknownSizes  int64                  `json:"unknownSizes"`
}
type RetryCounts struct {
	Selected int `json:"selected"`
	Queued   int `json:"queued"`
	Skipped  int `json:"skipped"`
}
type GCCounts struct {
	Visited        int   `json:"visited"`
	Reclaimed      int   `json:"reclaimed"`
	ReclaimedBytes int64 `json:"reclaimedBytes"`
	Done           bool  `json:"done"`
	Partial        bool  `json:"partial"`
}
type NativeTarget struct {
	MediaType string `json:"mediaType"`
	TMDBID    int64  `json:"tmdbId"`
	Season    *int   `json:"season,omitempty"`
	PageRequest
}
type NativeMetadata struct {
	State          string                `json:"state"`
	Name           string                `json:"name"`
	OriginalTitle  string                `json:"originalTitle"`
	Overview       string                `json:"overview"`
	RuntimeMinutes int                   `json:"runtimeMinutes"`
	Genres         []string              `json:"genres"`
	Episodes       Page[EpisodeMetadata] `json:"episodes"`
}
type EpisodeMetadata struct {
	Season         int    `json:"season"`
	Number         int    `json:"number"`
	Name           string `json:"name"`
	AirDate        string `json:"airDate"`
	RuntimeMinutes int    `json:"runtimeMinutes"`
}
