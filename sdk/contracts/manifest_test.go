package contracts_test

import (
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"os"
	"strings"
	"testing"
)

func example(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	b, e := os.ReadFile("../../plugins/hello-go/manifest.json")
	if e != nil {
		t.Fatal(e)
	}
	var m map[string]json.RawMessage
	if e = json.Unmarshal(b, &m); e != nil {
		t.Fatal(e)
	}
	return m
}
func marshal(t *testing.T, v any) []byte {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestManifestRequiredAndRetiredFields(t *testing.T) {
	valid := example(t)
	if _, e := contracts.ParseManifest(marshal(t, valid)); e != nil {
		t.Fatal(e)
	}
	for _, k := range []string{"contractVersion", "id", "name", "version", "author", "description", "level", "icon", "history", "permissions", "extensions", "ui", "resources", "configSchemaVersion", "dataSchemaVersion"} {
		t.Run("missing_"+k, func(t *testing.T) {
			m := example(t)
			delete(m, k)
			if _, e := contracts.ParseManifest(marshal(t, m)); e == nil {
				t.Fatal("required field accepted")
			}
		})
	}
	for _, k := range []string{"runtime", "sdkVersion", "entrypoint", "binarySHA256", "repository", "sourceURL", "downloadURL", "requirements", "packageDigest", "unknown"} {
		t.Run("retired_"+k, func(t *testing.T) {
			m := example(t)
			m[k] = json.RawMessage(`"blocked"`)
			if _, e := contracts.ParseManifest(marshal(t, m)); e == nil {
				t.Fatal("retired field accepted")
			}
		})
	}
	for _, tt := range []struct{ k, raw string }{{"id", `"../other"`}, {"icon", `"https://example.invalid/icon.svg"`}, {"icon", `"a/../icon.svg"`}, {"version", `"main"`}, {"level", `10`}, {"level", `null`}, {"permissions", `[{"scope":"events.*"}]`}, {"permissions", `[{"scope":"system.exec"},{"scope":"system.exec"}]`}, {"permissions", `[{"scope":"fs.data","resources":{"unknown":["x"]}}]`}, {"resources", `{"addressSpaceMiB":9000}`}, {"extensions", `{"api":[{"id":"x","path":"a","auth":"public","methods":["TRACE"],"executionScopes":[]}]}`}, {"ui", `{"mode":"vue","entry":"https://example.invalid/ui.js"}`}, {"ui", `{"mode":"vue","entry":"ui/dist/../escape.js"}`}} {
		t.Run(tt.k+tt.raw, func(t *testing.T) {
			m := example(t)
			m[tt.k] = json.RawMessage(tt.raw)
			if _, e := contracts.ParseManifest(marshal(t, m)); e == nil {
				t.Fatal("invalid metadata accepted")
			}
		})
	}
	for _, b := range []string{`{"id":"x","id":"y"}`, strings.Repeat(" ", 128<<10) + `{}`, `[]`} {
		if _, e := contracts.ParseManifest([]byte(b)); e == nil {
			t.Fatal("malformed manifest accepted")
		}
	}
}
func TestEventMapGoldenAndScopeValidation(t *testing.T) {
	b, e := os.ReadFile("testdata/event-map.json")
	if e != nil {
		t.Fatal(e)
	}
	var rows []struct{ Kind, Name, Wire, Type, State, Producer, Phase string }
	if e = json.Unmarshal(b, &rows); e != nil {
		t.Fatal(e)
	}
	counts := map[string]int{}
	na := 0
	seen := map[string]bool{}
	for _, r := range rows {
		key := r.Kind + ":" + r.Type
		if seen[key] {
			t.Fatal("duplicate event", key)
		}
		seen[key] = true
		counts[r.Kind]++
		known := false
		prefix := "events.subscribe."
		if r.Kind == "event" {
			known = contracts.KnownEventType(contracts.EventType(r.Type))
		} else {
			known = contracts.KnownChainEventType(contracts.ChainEventType(r.Type))
			prefix = "events.chain."
		}
		if r.State == "not_applicable" {
			na++
			if known {
				t.Fatal("N/A exposed")
			}
		} else if !known || !contracts.KnownScope(prefix+r.Type) {
			t.Fatal("missing known event/scope", key)
		}
	}
	if counts["event"] != 31 || counts["chain"] != 19 || na != 4 {
		t.Fatal(counts, na)
	}
	for _, s := range []string{"*", "events.*", "module.*", "events.subscribe.<type>", "module.unknown", "events.chain.agent.llm.provider", "settings.dumpEnv"} {
		if contracts.KnownScope(s) {
			t.Fatal("invalid scope", s)
		}
	}
	for _, s := range []string{"plugin.data.read", "system.exec", "events.subscribe.library.organized", "events.chain.name.recognize", "module.metadata.name_recognize"} {
		if !contracts.KnownScope(s) {
			t.Fatal("known scope missing", s)
		}
	}
}
func TestJSONAndErrors(t *testing.T) {
	for _, s := range []string{`{"x":1,"x":2}`, `{"x":[{"y":1,"y":2}]}`, `{} {}`, "\xef\xbb\xbf{}", `[NaN]`} {
		if contracts.ValidateJSON([]byte(s)) == nil {
			t.Fatal(s)
		}
	}
	for _, c := range []contracts.Code{contracts.ParseError, contracts.InvalidRequest, contracts.MethodNotFound, contracts.InvalidParams, contracts.InternalError, contracts.PermissionDenied, contracts.Unavailable, contracts.DeadlineExceeded, contracts.Cancelled, contracts.Conflict, contracts.Busy, contracts.NotFound, contracts.UnsupportedVersion, contracts.UpstreamLimited, contracts.NeedsAction, contracts.ResultUnknown, contracts.PayloadTooLarge, contracts.InvalidExtensionResult} {
		e := contracts.NewError(c)
		if contracts.Reason(c) != e.Data.Reason || strings.Contains(e.Error(), "http") {
			t.Fatal(e)
		}
	}
	e := contracts.SafeError(&contracts.Error{Code: 99, Message: "secret", Data: contracts.ErrorData{Reason: "secret", RetryAfterMS: 999999999999}})
	if e.Code != contracts.InternalError || e.Data.RetryAfterMS != 0 {
		t.Fatal(e)
	}
}
