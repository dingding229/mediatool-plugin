package contracts_test

import (
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestOrganizingContractsGoldenNarrowingAndPaths(t *testing.T) {
	q := contracts.TransferPlan{OperationID: aid, Paths: []string{"Movies/demo.strm"}}
	for _, p := range []string{"/tmp/a.strm", "../a.strm", "x/../../a.strm", "x//a.strm", "C:/a.strm", "x\\a.strm", ".media-tool-trash/a.strm", "a\n.strm"} {
		if contracts.ValidateTransferPaths(q, contracts.TransferPaths{Paths: []string{p}}) == nil {
			t.Fatal("unsafe path", p)
		}
	}
	if contracts.ValidateTransferPaths(q, contracts.TransferPaths{Paths: []string{"Movies/new.strm"}}) != nil {
		t.Fatal("safe rename refused")
	}
	if contracts.ValidateTransferPaths(q, contracts.TransferPaths{Paths: []string{"Movies/new.mp4"}}) == nil {
		t.Fatal("changed artifact kind")
	}
	if contracts.ValidateTransferPaths(q, contracts.TransferPaths{Paths: []string{}}) == nil {
		t.Fatal("coverage removed")
	}
	if contracts.ValidateDecision(contracts.TransferDecision{Decision: "allow_overwrite"}) == nil {
		t.Fatal("protection bypass")
	}
	storage := contracts.StorageSelection{Handles: []string{aid}, Selected: aid}
	if contracts.ValidateStorage(storage, contracts.StorageChoice{Handle: bid}) == nil {
		t.Fatal("unapproved root")
	}
	if contracts.ValidateStorage(storage, contracts.StorageChoice{Handle: aid}) != nil {
		t.Fatal("original handle refused")
	}
	for _, raw := range []string{`{"paths":[],"strm":"http://attacker"}`, `{"decision":"pass","force":true}`, `{"handle":"x","credentials":"private"}`} {
		var a contracts.TransferPaths
		var b contracts.TransferDecision
		var c contracts.StorageChoice
		if contracts.Decode([]byte(raw), &a) == nil || contracts.Decode([]byte(raw), &b) == nil || contracts.Decode([]byte(raw), &c) == nil {
			t.Fatal("unknown field", raw)
		}
	}
	plugintest.Golden(t, "testdata/organizing.json", map[string]any{
		"library.organize.failed":  contracts.OrganizeFailure{OperationID: aid, Kind: "transfer", ReasonCode: "attempts_exhausted"},
		"library.subtitles.copied": contracts.SubtitleOutcome{OperationID: aid, AssetID: bid, Count: 1, State: "copied"},
		"library.subtitles.failed": contracts.SubtitleOutcome{OperationID: aid, AssetID: bid, Count: 1, State: "failed"},
		"library.history.hidden":   contracts.HistoryHidden{RecordID: aid, Kind: "file"},
		"download.files.deleted":   contracts.DownloadDeleted{TaskID: aid, FilesDeleted: true},
		"download.task.deleted":    contracts.DownloadDeleted{TaskID: aid},
		"transfer.request":         q, "transfer.paths": contracts.TransferPaths{Paths: q.Paths}, "transfer.decision": contracts.TransferDecision{Decision: "pass"},
		"storage.request": storage, "storage.choice": contracts.StorageChoice{Handle: aid},
	})
	b, e := json.Marshal(q)
	if e != nil || contracts.ValidateJSON(b) != nil {
		t.Fatal(e)
	}
}
