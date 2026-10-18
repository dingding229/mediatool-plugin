package contracts

import (
	"path"
	"strings"
	"unicode/utf8"
)

type OrganizeFailure struct {
	OperationID string `json:"operationId"`
	Kind        string `json:"kind"`
	ReasonCode  string `json:"reasonCode"`
}
type SubtitleOutcome struct {
	OperationID string `json:"operationId"`
	AssetID     string `json:"assetId"`
	Count       int    `json:"count"`
	State       string `json:"state"`
}
type HistoryHidden struct {
	RecordID          string `json:"recordId"`
	Kind              string `json:"kind"`
	FilesystemDeleted bool   `json:"filesystemDeleted"`
}
type DownloadDeleted struct {
	TaskID       string `json:"taskId"`
	FilesDeleted bool   `json:"filesDeleted"`
}

// Paths are relative proposals only. No root, resolver URL, source identity,
// playback grant or STRM bytes can be represented by a rename response.
type TransferPlan struct {
	OperationID string   `json:"operationId"`
	Paths       []string `json:"paths"`
}
type TransferPaths struct {
	Paths []string `json:"paths"`
}
type TransferDecision struct {
	Decision string `json:"decision"`
}
type StorageSelection struct {
	Handles  []string `json:"handles"`
	Selected string   `json:"selected"`
}
type StorageChoice struct {
	Handle string `json:"handle"`
}

func RelativeOutput(p string) bool {
	if p == "" || len(p) > 2048 || !utf8.ValidString(p) || path.IsAbs(p) || path.Clean(p) != p || p == "." || strings.ContainsAny(p, "\\:\x00\r\n") {
		return false
	}
	for _, part := range strings.Split(p, "/") {
		if part == ".." || part == "." || part == "" || strings.HasPrefix(part, ".media-tool-") || strings.TrimSpace(part) != part {
			return false
		}
	}
	return true
}
func ValidateTransfer(v TransferPlan) error {
	if !OpaqueID(v.OperationID) || len(v.Paths) > 100 {
		return NewError(InvalidParams)
	}
	seen := map[string]bool{}
	for _, p := range v.Paths {
		if !RelativeOutput(p) || seen[p] {
			return NewError(InvalidParams)
		}
		seen[p] = true
	}
	return nil
}
func ValidateTransferPaths(v TransferPlan, r TransferPaths) error {
	if ValidateTransfer(v) != nil || r.Paths == nil || len(r.Paths) != len(v.Paths) {
		return NewError(InvalidExtensionResult)
	}
	seen := map[string]bool{}
	for i, p := range r.Paths {
		if !RelativeOutput(p) || seen[p] || path.Ext(p) != path.Ext(v.Paths[i]) {
			return NewError(InvalidExtensionResult)
		}
		seen[p] = true
	}
	return nil
}
func ValidateDecision(r TransferDecision) error {
	if r.Decision != "pass" && r.Decision != "reject" && r.Decision != "skip" {
		return NewError(InvalidExtensionResult)
	}
	return nil
}
func ValidateStorage(v StorageSelection, r StorageChoice) error {
	if len(v.Handles) < 1 || len(v.Handles) > 200 || r.Handle == "" {
		return NewError(InvalidExtensionResult)
	}
	seen := map[string]bool{}
	allowed := false
	for _, h := range v.Handles {
		if !OpaqueID(h) || seen[h] {
			return NewError(InvalidExtensionResult)
		}
		seen[h] = true
		allowed = allowed || h == r.Handle
	}
	if !seen[v.Selected] || !allowed {
		return NewError(InvalidExtensionResult)
	}
	return nil
}
