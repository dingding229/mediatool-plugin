package business

import "time"

type Empty struct{}
type Text struct {
	Text string `json:"text"`
}
type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Asset struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	MediaType string `json:"mediaType"`
	Year      string `json:"year"`
	State     string `json:"state"`
	Revision  int64  `json:"revision"`
	Revoked   bool   `json:"revoked"`
}
type Candidate struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	AssetID  string `json:"assetId,omitempty"`
	State    string `json:"state"`
	Revision int64  `json:"revision"`
	Size     int64  `json:"size"`
}
type LibraryPage struct {
	PageRequest
	State string `json:"state,omitempty"`
	Kind  string `json:"kind,omitempty"`
}
type Scan struct {
	ID          string `json:"id"`
	RootID      string `json:"rootId"`
	State       string `json:"state"`
	Files       int    `json:"files"`
	Directories int    `json:"directories"`
	Revision    int64  `json:"revision"`
}
type LibraryHistory struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	State     string    `json:"state"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type ScanStart struct {
	RootID string `json:"rootId"`
	Path   string `json:"path"`
}
type Control struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}
type Recognition struct {
	Candidate         Candidate `json:"candidate"`
	Hint              NameHint  `json:"hint"`
	Matches           []Media   `json:"matches"`
	Selected          *Media    `json:"selected,omitempty"`
	FileDigest        string    `json:"fileDigest"`
	NeedsConfirmation bool      `json:"needsConfirmation"`
}
type ConfirmIdentity struct {
	ID             string `json:"id"`
	MediaType      string `json:"mediaType"`
	TMDBID         int    `json:"tmdbId"`
	Season         int    `json:"season"`
	Episodes       []int  `json:"episodes"`
	Revision       int64  `json:"revision"`
	MigrateOutputs bool   `json:"migrateOutputs"`
	FileDigest     string `json:"fileDigest"`
	PlanDigest     string `json:"planDigest"`
}
type IdentityPlan struct {
	ID         string `json:"id"`
	PlanDigest string `json:"planDigest"`
	Revision   int64  `json:"revision"`
	Conflicts  int    `json:"conflicts"`
}
type HistorySelection struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type HistorySelections struct {
	Items []HistorySelection `json:"items"`
}
type Archive struct {
	ID            string    `json:"id"`
	AssetID       string    `json:"assetId"`
	Kind          string    `json:"kind"`
	State         string    `json:"state"`
	EligibleAfter time.Time `json:"eligibleAfter"`
}
type ArchivePlan struct {
	Items            []Archive `json:"items"`
	PlanDigest       string    `json:"planDigest"`
	PhysicalDeletion bool      `json:"physicalDeletion"`
}
type ArchiveApproval struct {
	IDs        []string `json:"ids"`
	PlanDigest string   `json:"planDigest"`
}
type CleanupResult struct {
	ID        string `json:"id"`
	ArchiveID string `json:"archiveId"`
	State     string `json:"state"`
}
type Root struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	Revision int64  `json:"revision"`
	Read     bool   `json:"read"`
	Write    bool   `json:"write"`
}
type FileQuery struct {
	RootID string `json:"rootId"`
	Path   string `json:"path"`
	PageRequest
}
type MediaFile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Size        int64  `json:"size"`
	Directory   bool   `json:"directory"`
	Unsupported bool   `json:"unsupported"`
}
type FilePage struct {
	Page[MediaFile]
	RootID string `json:"rootId"`
}

// FileInput contains only root-relative names. Preview returns an expiring evidence digest;
// Apply requires that exact digest and an explicit confirmed Operation. No raw upload/copy.
type FileInput struct {
	RootID         string `json:"rootId"`
	Operation      string `json:"operation"`
	Source         string `json:"source"`
	Destination    string `json:"destination"`
	Name           string `json:"name"`
	ExpectedFileID string `json:"expectedFileId"`
	ExpectedSize   *int64 `json:"expectedSize,omitempty"`
}
type FileInputs struct {
	Items []FileInput `json:"items"`
}
type Approval struct {
	ID     string `json:"id"`
	Digest string `json:"digest"`
}
type FilePlan struct {
	ID             string      `json:"id"`
	RootID         string      `json:"rootId"`
	Revision       int64       `json:"revision"`
	Operation      string      `json:"operation"`
	State          string      `json:"state"`
	Digest         string      `json:"digest"`
	ExpiresAt      time.Time   `json:"expiresAt"`
	Items          []MediaFile `json:"items"`
	CompletedSteps int         `json:"completedSteps"`
}
type FileBatch struct {
	ID     string     `json:"id"`
	Digest string     `json:"digest"`
	State  string     `json:"state"`
	Items  []FilePlan `json:"items"`
}
type RepairItem struct {
	AssetID  string `json:"assetId"`
	Revision int64  `json:"revision"`
	State    string `json:"state"`
	Digest   string `json:"digest"`
}
type RepairPlan struct {
	ID        string       `json:"id"`
	Digest    string       `json:"digest"`
	State     string       `json:"state"`
	ExpiresAt time.Time    `json:"expiresAt"`
	Items     []RepairItem `json:"items"`
}
type RenderName struct {
	Template string   `json:"template"`
	Hint     NameHint `json:"hint"`
	Original string   `json:"original"`
	Path     bool     `json:"path"`
}
type EpisodeCoverage struct {
	Season   int   `json:"season"`
	Episodes []int `json:"episodes"`
	Valid    bool  `json:"valid"`
}
type MetadataQuery struct {
	MediaType string `json:"mediaType"`
	ID        int    `json:"id"`
	Query     string `json:"query"`
	Season    *int   `json:"season,omitempty"`
}
type StoredMetadata struct {
	State          string `json:"state"`
	Media          Media  `json:"media"`
	Overview       string `json:"overview"`
	RuntimeMinutes int    `json:"runtimeMinutes"`
}
type Classify struct {
	Media            Media    `json:"media"`
	GenreIDs         []int    `json:"genreIds"`
	Countries        []string `json:"countries"`
	OriginalLanguage string   `json:"originalLanguage"`
	Source           string   `json:"source"`
	MovieTemplate    string   `json:"movieTemplate"`
	SeriesTemplate   string   `json:"seriesTemplate"`
	Extension        string   `json:"extension"`
	Season           int      `json:"season"`
	Episodes         []int    `json:"episodes"`
}
type NameHint struct {
	Title      string   `json:"title"`
	Year       string   `json:"year"`
	MediaType  string   `json:"mediaType"`
	Season     int      `json:"season"`
	Episode    int      `json:"episode"`
	Episodes   []int    `json:"episodes"`
	Resolution string   `json:"resolution"`
	Source     string   `json:"source"`
	Codec      string   `json:"codec"`
	Group      string   `json:"group"`
	Extension  string   `json:"extension"`
	Confidence string   `json:"confidence"`
	Warnings   []string `json:"warnings"`
}
type ClassRule struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Directory string   `json:"directory"`
	Enabled   bool     `json:"enabled"`
	GenreIDs  []int    `json:"genreIds"`
	Countries []string `json:"countries"`
	Languages []string `json:"languages"`
}
type ClassGroup struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Directory string      `json:"directory"`
	MediaType string      `json:"mediaType"`
	GenreIDs  []int       `json:"genreIds"`
	Enabled   bool        `json:"enabled"`
	Rules     []ClassRule `json:"rules"`
}
type ClassPolicy struct {
	Groups []ClassGroup `json:"groups"`
}
type Classification struct {
	Value    ClassPolicy `json:"value"`
	Revision int64       `json:"revision"`
}
type Category struct {
	GroupID   string `json:"groupId"`
	RuleID    string `json:"ruleId"`
	Name      string `json:"name"`
	Directory string `json:"directory"`
	Reason    string `json:"reason"`
}
