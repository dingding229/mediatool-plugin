package contracts_test

import (
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestDiscoveryGoldenAndNarrowing(t *testing.T) {
	c := contracts.MediaIdentity{TMDBID: 1, MediaType: "series"}
	q := contracts.MediaConversion{Candidates: []contracts.MediaIdentity{c}, Selected: &c}
	if contracts.ValidateConversion(q, contracts.MediaChoice{Selected: &c}) != nil {
		t.Fatal("legal choice")
	}
	for _, v := range []contracts.MediaIdentity{{TMDBID: 2, MediaType: "series"}, {TMDBID: 1, MediaType: "movie"}, {TMDBID: -1, MediaType: "tv"}} {
		if contracts.ValidateConversion(q, contracts.MediaChoice{Selected: &v}) == nil {
			t.Fatal("unknown identity")
		}
	}
	e := contracts.EpisodeExpectation{SubscriptionID: aid, Revision: 1, TMDBID: 1, Season: 1, Expected: []int{1, 2}}
	if contracts.ValidateExpectation(e, contracts.EpisodeSupplement{Expected: []int{1}}) == nil {
		t.Fatal("removed expectation")
	}
	if contracts.ValidateExpectation(e, contracts.EpisodeSupplement{Expected: []int{1, 2, 3}}) != nil {
		t.Fatal("supplement")
	}
	check := contracts.CompletionCheck{EpisodeExpectation: e, Owned: []int{1}, Complete: false}
	if contracts.ValidateCompletion(check, contracts.CompletionDecision{Allow: true}) == nil {
		t.Fatal("invented ownership")
	}
	for _, title := range []string{"http://private/token", "/data/private", "Bearer:secret", "a\nb"} {
		if contracts.ValidateSource(contracts.SourceRequest{Page: 1, Size: 200}, contracts.SourceResult{Items: []contracts.SourceCandidate{{MediaIdentity: c, Title: title}}}) == nil {
			t.Fatal("unsafe source")
		}
	}
	completed := check
	completed.Owned = []int{1, 2}
	completed.Complete = true
	plugintest.Golden(t, "testdata/discovery.json", map[string]any{"metadata.scraped": contracts.MetadataScraped{ObjectID: aid, Kind: "asset", MediaIdentity: c, Revision: "v1"}, "subscription.coverage.completed": contracts.CoverageObservation{CompletionCheck: completed, State: "complete", Evidence: "exposable_assets", ObservedAt: "2026-10-18T00:00:00Z"}, "media.recognize.convert": q, "subscribe.episodes.refresh": e, "subscribe.completion.check": check, "discover.source": contracts.SourceResult{Items: []contracts.SourceCandidate{{MediaIdentity: c, Title: "Demo", Year: 2026}}}})
}
