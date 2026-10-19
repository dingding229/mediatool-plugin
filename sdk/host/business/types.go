// Package business defines stable, credential-free acquisition-chain DTOs. ContractVersion remains 2.
package business

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/host"
	"strings"
	"unicode/utf8"
)

type PageRequest struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Query string `json:"query,omitempty"`
}

func (p PageRequest) Valid() bool {
	return p.Page >= 1 && p.Page <= 1000 && p.Size >= 1 && p.Size <= 200 && utf8.ValidString(p.Query) && utf8.RuneCountInString(p.Query) <= 200 && !strings.ContainsAny(p.Query, "\x00\r\n")
}

type Page[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"hasMore"`
}
type ID struct {
	ID string `json:"id"`
}
type Selection struct {
	ID string `json:"id"`
	PageRequest
}
type Operation struct {
	Key      string
	Revision *int64
	Confirm  bool
}
type Receipt struct {
	ID     string `json:"id,omitempty"`
	State  string `json:"state"`
	Cached bool   `json:"cached,omitempty"`
}
type Media struct {
	ID            int    `json:"id"`
	MediaType     string `json:"mediaType"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle,omitempty"`
	Year          string `json:"year,omitempty"`
}
type Resource struct {
	ID             string   `json:"id"`
	AccountID      string   `json:"accountId"`
	Title          string   `json:"title"`
	Size           int64    `json:"size"`
	Seeders        int      `json:"seeders"`
	Leechers       int      `json:"leechers"`
	HR             *bool    `json:"hr,omitempty"`
	DownloadFactor *float64 `json:"downloadFactor,omitempty"`
}
type Task struct {
	SeedingState    string     `json:"seedingState,omitempty"`
	SiteID          string     `json:"siteId,omitempty"`
	Category        string     `json:"category,omitempty"`
	Tags            []string   `json:"tags,omitempty"`
	Ratio           float64    `json:"ratio"`
	SeedingSeconds  int64      `json:"seedingSeconds"`
	Complete        bool       `json:"complete"`
	Seeding         bool       `json:"seeding"`
	CleanupEligible bool       `json:"cleanupEligible"`
	Transfers       []Transfer `json:"transfers,omitempty"`
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	State           string     `json:"state"`
	Hash            string     `json:"hash"`
	Owned           bool       `json:"owned"`
	Progress        *float64   `json:"progress,omitempty"`
	DownloadSpeed   int64      `json:"downloadSpeed"`
	UploadSpeed     int64      `json:"uploadSpeed"`
	Size            int64      `json:"size"`
	Revision        int64      `json:"revision,omitempty"`
	CleanupState    string     `json:"cleanupState,omitempty"`
}
type File struct {
	Index    int     `json:"index"`
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Progress float64 `json:"progress"`
}
type Tracker struct {
	Origin string `json:"origin"`
	Status int    `json:"status"`
	Seeds  int    `json:"seeds"`
	Peers  int    `json:"peers"`
}
type AddResource struct {
	RunID        string `json:"runId"`
	ResourceID   string `json:"resourceId"`
	MediaType    string `json:"mediaType"`
	TMDBID       int    `json:"tmdbId"`
	ExpectedHash string `json:"expectedHash"`
	Season       *int   `json:"season,omitempty"`
	Episodes     []int  `json:"episodes,omitempty"`
}
type ManualSource struct {
	Magnet  string `json:"magnet,omitempty"`
	Torrent []byte `json:"torrent,omitempty"`
}
type Retry struct {
	ID                  string `json:"id"`
	AcceptDuplicateRisk bool   `json:"acceptDuplicateRisk"`
}
type Paused struct {
	Hash   string `json:"hash"`
	Paused bool   `json:"paused"`
}
type Cleanup struct {
	ID           string `json:"id"`
	Hash         string `json:"hash"`
	Eligible     bool   `json:"eligible"`
	State        string `json:"state"`
	DeletesFiles bool   `json:"deletesFiles"`
}
type Torrent struct {
	Hash    string   `json:"hash"`
	Aliases []string `json:"aliases"`
	Version int      `json:"version"`
	Name    string   `json:"name"`
	Files   []File   `json:"files"`
}
type TorrentInput struct {
	Data []byte `json:"data"`
}
type SearchStart struct {
	Keyword    string   `json:"keyword"`
	AccountIDs []string `json:"accountIds"`
	Refresh    bool     `json:"refresh"`
	Media      *Media   `json:"media,omitempty"`
}
type SearchRun struct {
	ID        string     `json:"id"`
	Keyword   string     `json:"keyword"`
	State     string     `json:"state"`
	Resources []Resource `json:"resources"`
	HasMore   bool       `json:"hasMore"`
}
type Site struct {
	Origin     string `json:"origin"`
	TemplateID string `json:"templateId"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	Enabled    bool   `json:"enabled"`
	Revision   int64  `json:"revision"`
}
type Metrics struct {
	AccountID    string   `json:"accountId"`
	State        string   `json:"state"`
	Uploaded     *int64   `json:"uploaded,omitempty"`
	Downloaded   *int64   `json:"downloaded,omitempty"`
	Ratio        *float64 `json:"ratio,omitempty"`
	Bonus        *float64 `json:"bonus,omitempty"`
	Seeding      *int64   `json:"seeding,omitempty"`
	SeedingBytes *int64   `json:"seedingBytes,omitempty"`
}
type SiteRequest struct {
	AccountID     string `json:"accountId"`
	Origin        string `json:"origin"`
	Path          string `json:"path"`
	Method        string `json:"method"`
	Body          string `json:"body,omitempty"`
	SuccessMarker string `json:"successMarker,omitempty"`
	Raw           bool   `json:"raw"`
}
type SiteResponse struct {
	Status  int    `json:"status"`
	State   string `json:"state"`
	Matched bool   `json:"matched"`
	Body    string `json:"body,omitempty"`
}
type Template struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Loaded bool   `json:"loaded"`
}
type IDs struct {
	IDs []string `json:"ids"`
}
type Rules struct {
	Resolution     string   `json:"resolution,omitempty"`
	Source         string   `json:"source,omitempty"`
	Codec          string   `json:"codec,omitempty"`
	AccountIDs     []string `json:"accountIds"`
	RegexKeywords  bool     `json:"regexKeywords,omitempty"`
	Keyword        string   `json:"keyword"`
	Include        []string `json:"include"`
	Exclude        []string `json:"exclude"`
	MinSize        int64    `json:"minSize"`
	MaxSize        int64    `json:"maxSize"`
	MinSeeders     int      `json:"minSeeders"`
	RejectHR       *bool    `json:"rejectHR,omitempty"`
	FreeOnly       bool     `json:"freeOnly"`
	Season         int      `json:"season"`
	Episodes       []int    `json:"episodes"`
	IncludeHistory bool     `json:"includeHistory"`
}
type Subscription struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Media    Media  `json:"media"`
	Enabled  bool   `json:"enabled"`
	Source   string `json:"source"`
	Rules    Rules  `json:"rules"`
	Revision int64  `json:"revision"`
	State    string `json:"state"`
}
type SubscriptionInput struct {
	ID        string `json:"id,omitempty"`
	MediaType string `json:"mediaType"`
	TMDBID    int    `json:"tmdbId"`
	Enabled   bool   `json:"enabled"`
	Source    string `json:"source"`
	Rules     Rules  `json:"rules"`
}
type Queue struct {
	ID      string `json:"id"`
	Preview bool   `json:"preview"`
}
type History struct {
	ID      string `json:"id"`
	State   string `json:"state"`
	Preview bool   `json:"preview"`
}
type Evaluation struct {
	Resource Resource `json:"resource"`
	Media    Media    `json:"media"`
	Rules    Rules    `json:"rules"`
}
type Evaluated struct {
	Accepted bool     `json:"accepted"`
	Reasons  []string `json:"reasons"`
}
type CoverageInput struct {
	Names  []string `json:"names"`
	Season int      `json:"season"`
}
type Coverage struct {
	Episodes []int  `json:"episodes"`
	Evidence string `json:"evidence"`
}

func Call[T any](ctx context.Context, c *host.Client, method string, args any, op Operation) (T, error) {
	var out T
	s, ok := Lookup(method)
	if !ok {
		return out, contracts.NewError(contracts.MethodNotFound)
	}
	if s.Level != "R" && (!op.Confirm || len(op.Key) < 8 || len(op.Key) > 100) {
		return out, contracts.NewError(contracts.InvalidParams)
	}
	e := c.Invoke(ctx, method, args, op.Key, op.Revision, &out)
	return out, e
}

// ValidateArgs is deliberately repeated by the host; SDK checks are not an authorization boundary.
func ValidateArgs(raw json.RawMessage) error {
	if len(raw) > 128<<10 {
		return contracts.NewError(contracts.PayloadTooLarge)
	}
	return contracts.ValidateJSON(raw)
}

// Upload controls address a program-owned transfer ID, not a qB download ID.
type UploadControl struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}
type Transfer struct {
	ID       string `json:"id"`
	State    string `json:"state"`
	Progress int    `json:"progress"`
}

// DeleteFiles is always explicit. False retires only the owned, imported,
// HR-cleared torrent; true requests native cleanup and all its file/cloud proofs.
type DeleteRequest struct {
	ID          string `json:"id"`
	DeleteFiles bool   `json:"deleteFiles"`
}
