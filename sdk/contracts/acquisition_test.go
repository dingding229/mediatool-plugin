package contracts_test

import (
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
	"time"
)

const aid = "00000000-0000-0000-0000-000000000001"
const bid = "00000000-0000-0000-0000-000000000002"

func TestAcquisitionContractsGoldenAndNarrowing(t *testing.T) {
	q := contracts.ResourceSelectionRequest{RunID: aid, Candidates: []contracts.ResourceCandidate{{ID: aid, AccountID: bid, Size: 100, Seeders: 3}, {ID: bid, AccountID: bid}}}
	for _, ids := range [][]string{{bid, aid}, {bid}, {}} {
		if e := contracts.ValidateSelection(q, contracts.ResourceSelectionResult{IDs: ids}); e != nil {
			t.Fatal(e)
		}
	}
	for _, ids := range [][]string{nil, {aid, aid}, {"00000000-0000-0000-0000-000000000003"}} {
		if contracts.ValidateSelection(q, contracts.ResourceSelectionResult{IDs: ids}) == nil {
			t.Fatal("expanded candidate set", ids)
		}
	}
	narrowed := q
	narrowed.Candidates = narrowed.Candidates[:1]
	if contracts.ValidateSelection(narrowed, contracts.ResourceSelectionResult{IDs: []string{bid}}) == nil {
		t.Fatal("later handler restored removed ID")
	}
	request := contracts.ResourceDownloadRequest{Kind: "managed", ResourceID: aid, DownloaderID: bid, CleanupAfterImport: true}
	if contracts.ValidateDownload(request, contracts.ResourceDownloadResult{Allow: true}) != nil {
		t.Fatal("cleanup removal refused")
	}
	request.CleanupAfterImport = false
	if contracts.ValidateDownload(request, contracts.ResourceDownloadResult{Allow: true, CleanupAfterImport: true}) == nil {
		t.Fatal("cleanup permission widened")
	}
	for _, raw := range []string{`{"ids":[],"url":"https://attacker.invalid"}`, `{"allow":true,"cleanupAfterImport":false,"ignoreHR":true}`, `{"allow":true,"path":"/private"}`} {
		var v contracts.ResourceSelectionResult
		var d contracts.ResourceDownloadResult
		if contracts.Decode([]byte(raw), &v) == nil || contracts.Decode([]byte(raw), &d) == nil {
			t.Fatal("unknown fields accepted")
		}
	}
	plugintest.Golden(t, "testdata/acquisition.json", map[string]any{
		"site.account.updated":       contracts.SiteAccountChanged{AccountID: aid, Revision: 2, Enabled: true},
		"site.account.deleted":       contracts.SiteAccountChanged{AccountID: aid, Revision: 2},
		"site.metrics.refreshed":     contracts.SiteMetricsRefreshed{AccountID: aid, AccountRevision: 2, State: "succeeded", AttemptedAt: time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)},
		"subscription.created":       contracts.SubscriptionChanged{SubscriptionID: bid, Revision: 1, Enabled: true, MediaType: "movie", TMDBID: 1},
		"subscription.updated":       contracts.SubscriptionChanged{SubscriptionID: bid, Revision: 2, MediaType: "movie", TMDBID: 1},
		"subscription.deleted":       contracts.SubscriptionChanged{SubscriptionID: bid, Revision: 3, MediaType: "movie", TMDBID: 1},
		"resource.selection.request": q, "resource.selection.result": contracts.ResourceSelectionResult{IDs: []string{bid, aid}},
		"resource.download.request": request, "resource.download.result": contracts.ResourceDownloadResult{Allow: true},
	})
	b, _ := json.Marshal(q)
	if contracts.ValidateJSON(b) != nil {
		t.Fatal("invalid JSON")
	}
}
