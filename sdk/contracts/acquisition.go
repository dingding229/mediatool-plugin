package contracts

import (
	"strings"
	"time"
)

// Acquisition events intentionally carry identities/status only. No account
// settings, titles, credentials, tracker/download URLs or source hints cross RPC.
type SiteAccountChanged struct {
	AccountID string `json:"accountId"`
	Revision  int64  `json:"revision"`
	Enabled   bool   `json:"enabled"`
}
type SiteMetricsRefreshed struct {
	AccountID       string    `json:"accountId"`
	AccountRevision int64     `json:"accountRevision"`
	State           string    `json:"state"`
	AttemptedAt     time.Time `json:"attemptedAt"`
}
type SubscriptionChanged struct {
	SubscriptionID string `json:"subscriptionId"`
	Revision       int64  `json:"revision"`
	Enabled        bool   `json:"enabled"`
	MediaType      string `json:"mediaType"`
	TMDBID         int    `json:"tmdbId"`
}

// ResourceSelection can only reorder/remove opaque IDs from this frozen window
// (at most 200). Later windows cannot be moved ahead of this one.
type ResourceCandidate struct {
	ID        string `json:"id"`
	AccountID string `json:"accountId"`
	Size      int64  `json:"size"`
	Seeders   int    `json:"seeders"`
}
type ResourceSelectionRequest struct {
	RunID      string              `json:"runId"`
	Candidates []ResourceCandidate `json:"candidates"`
}
type ResourceSelectionResult struct {
	IDs []string `json:"ids"`
}

// ResourceDownload is policy only, never a downloader replacement. The only
// mutable option is removal of post-import cleanup; it cannot enable deletion.
// Manual downloads have cleanupAfterImport=false. No path/source/hash is mutable.
type ResourceDownloadRequest struct {
	Kind               string `json:"kind"`
	ResourceID         string `json:"resourceId"`
	DownloaderID       string `json:"downloaderId"`
	CleanupAfterImport bool   `json:"cleanupAfterImport"`
}
type ResourceDownloadResult struct {
	Allow              bool `json:"allow"`
	CleanupAfterImport bool `json:"cleanupAfterImport"`
}

func OpaqueID(s string) bool {
	if len(s) != 36 {
		return false
	}
	for i, c := range s {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !strings.ContainsRune("0123456789abcdef", c) {
			return false
		}
	}
	return true
}
func ValidateSelectionRequest(v ResourceSelectionRequest) error {
	if !OpaqueID(v.RunID) || len(v.Candidates) > 200 {
		return NewError(InvalidParams)
	}
	seen := map[string]bool{}
	for _, c := range v.Candidates {
		if !OpaqueID(c.ID) || !OpaqueID(c.AccountID) || seen[c.ID] || c.Size < 0 || c.Seeders < 0 {
			return NewError(InvalidParams)
		}
		seen[c.ID] = true
	}
	return nil
}
func ValidateSelection(v ResourceSelectionRequest, r ResourceSelectionResult) error {
	if ValidateSelectionRequest(v) != nil || r.IDs == nil || len(r.IDs) > len(v.Candidates) {
		return NewError(InvalidExtensionResult)
	}
	allowed := map[string]bool{}
	for _, c := range v.Candidates {
		allowed[c.ID] = true
	}
	for _, id := range r.IDs {
		if !allowed[id] {
			return NewError(InvalidExtensionResult)
		}
		delete(allowed, id)
	}
	return nil
}
func ValidateDownloadRequest(v ResourceDownloadRequest) error {
	if (v.Kind != "managed" && v.Kind != "manual") || !OpaqueID(v.ResourceID) || !OpaqueID(v.DownloaderID) || v.Kind == "manual" && v.CleanupAfterImport {
		return NewError(InvalidParams)
	}
	return nil
}
func ValidateDownload(v ResourceDownloadRequest, r ResourceDownloadResult) error {
	if ValidateDownloadRequest(v) != nil || !v.CleanupAfterImport && r.CleanupAfterImport {
		return NewError(InvalidExtensionResult)
	}
	return nil
}
