package contracts

import (
	"slices"
	"strings"
	"unicode"
)

// Identity-only contracts deliberately cannot represent artwork, credentials,
// player URLs, source paths or a claim that a file is owned.
type MediaIdentity struct {
	TMDBID    int    `json:"tmdbId"`
	MediaType string `json:"mediaType"`
}
type MediaConversion struct {
	Candidates []MediaIdentity `json:"candidates"`
	Selected   *MediaIdentity  `json:"selected"`
}
type MediaChoice struct {
	Selected *MediaIdentity `json:"selected"`
}
type MetadataScraped struct {
	ObjectID string `json:"objectId"`
	Kind     string `json:"kind"`
	MediaIdentity
	Revision string `json:"revision"`
}
type EpisodeExpectation struct {
	SubscriptionID string `json:"subscriptionId"`
	Revision       int64  `json:"revision"`
	TMDBID         int    `json:"tmdbId"`
	Season         int    `json:"season"`
	Expected       []int  `json:"expected"`
}
type EpisodeSupplement struct {
	Expected []int `json:"expected"`
}
type CompletionCheck struct {
	EpisodeExpectation
	Owned    []int `json:"owned"`
	Complete bool  `json:"complete"`
}
type CompletionDecision struct {
	Allow bool `json:"allow"`
}
type CoverageObservation struct {
	CompletionCheck
	State      string `json:"state"`
	Evidence   string `json:"evidence"`
	ObservedAt string `json:"observedAt"`
}
type SourceRequest struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
type SourceCandidate struct {
	MediaIdentity
	Title string `json:"title"`
	Year  int    `json:"year"`
}
type SourceResult struct {
	Items   []SourceCandidate `json:"items"`
	HasMore bool              `json:"hasMore"`
}

func ValidIdentity(v MediaIdentity) bool {
	return v.TMDBID > 0 && v.TMDBID <= 2147483647 && (v.MediaType == "movie" || v.MediaType == "series")
}
func ValidateConversion(v MediaConversion, r MediaChoice) error {
	if len(v.Candidates) > 200 {
		return NewError(InvalidExtensionResult)
	}
	for _, c := range v.Candidates {
		if !ValidIdentity(c) {
			return NewError(InvalidExtensionResult)
		}
	}
	if r.Selected != nil && (!ValidIdentity(*r.Selected) || !slices.Contains(v.Candidates, *r.Selected)) {
		return NewError(InvalidExtensionResult)
	}
	return nil
}
func ValidEpisodes(v []int) bool {
	if v == nil || len(v) > 200 {
		return false
	}
	seen := map[int]bool{}
	for _, n := range v {
		if n < 1 || n > 9999 || seen[n] {
			return false
		}
		seen[n] = true
	}
	return true
}
func ValidateExpectation(v EpisodeExpectation, r EpisodeSupplement) error {
	if !OpaqueID(v.SubscriptionID) || v.Revision < 1 || v.TMDBID < 1 || v.Season < 0 || v.Season > 999 || !ValidEpisodes(v.Expected) || !ValidEpisodes(r.Expected) {
		return NewError(InvalidExtensionResult)
	}
	for _, n := range v.Expected {
		if !slices.Contains(r.Expected, n) {
			return NewError(InvalidExtensionResult)
		}
	}
	return nil
}
func ValidateCompletion(v CompletionCheck, r CompletionDecision) error {
	if ValidateExpectation(v.EpisodeExpectation, EpisodeSupplement{Expected: v.Expected}) != nil || !ValidEpisodes(v.Owned) {
		return NewError(InvalidExtensionResult)
	}
	complete := len(v.Expected) > 0
	for _, n := range v.Expected {
		complete = complete && slices.Contains(v.Owned, n)
	}
	if r.Allow && (!v.Complete || !complete) {
		return NewError(InvalidExtensionResult)
	}
	return nil
}
func ValidateSource(q SourceRequest, v SourceResult) error {
	if q.Page < 1 || q.Page > 1000 || q.Size < 1 || q.Size > 200 || v.Items == nil || len(v.Items) > q.Size {
		return NewError(InvalidExtensionResult)
	}
	seen := map[MediaIdentity]bool{}
	for _, c := range v.Items {
		if !ValidIdentity(c.MediaIdentity) || seen[c.MediaIdentity] || !SafeSourceTitle(c.Title) || c.Year < 0 || c.Year > 9999 {
			return NewError(InvalidExtensionResult)
		}
		seen[c.MediaIdentity] = true
	}
	return nil
}
func SafeSourceTitle(s string) bool {
	if len([]rune(s)) < 1 || len([]rune(s)) > 200 || strings.TrimSpace(s) != s || strings.ContainsAny(s, "/\\@") || strings.Contains(s, ":") {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}
