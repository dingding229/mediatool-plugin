package business

// WishlistRequest selects a fixed public provider, never an arbitrary URL or credential.
type WishlistRequest struct {
	Source string `json:"source"`
	ID     string `json:"id"`
	Page   int    `json:"page"`
}
type WishlistEntry struct {
	SourceID  string `json:"sourceId"`
	Title     string `json:"title"`
	Year      string `json:"year"`
	MediaType string `json:"mediaType"`
	TMDBID    int    `json:"tmdbId"`
}
type WishlistPage struct {
	Items             []WishlistEntry `json:"items"`
	HasMore           bool            `json:"hasMore"`
	State             string          `json:"state"`
	RetryAfterSeconds int             `json:"retryAfterSeconds,omitempty"`
}

// Positive optional limits; zero/unlimited resets, save paths and credentials are not exposed.
type DownloaderConfig struct {
	ID             string   `json:"id"`
	Tag            string   `json:"tag,omitempty"`
	Category       string   `json:"category,omitempty"`
	DownloadLimit  *int64   `json:"downloadLimit,omitempty"`
	UploadLimit    *int64   `json:"uploadLimit,omitempty"`
	RatioLimit     *float64 `json:"ratioLimit,omitempty"`
	SeedingMinutes *int64   `json:"seedingMinutes,omitempty"`
}
