package registry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCatalogPureIndependentSnapshots(t *testing.T) {
	a, e := All()
	if e != nil || len(a) != 5 || a[0].Manifest.ID != "hello-go" {
		t.Fatal(e, a)
	}
	a[0].Manifest.Name = "changed"
	a[0].Manifest.History["1.0.0"] = "changed"
	a[0].Manifest.Permissions[0].Scope = "system.exec"
	b, e := All()
	if e != nil || b[0].Manifest.Name == "changed" || b[0].Manifest.History["1.0.0"] == "changed" || b[0].Manifest.Permissions[0].Scope == "system.exec" {
		t.Fatal("shared mutable metadata")
	}
	if _, e = Get("missing"); e == nil {
		t.Fatal("unknown plugin")
	}
	files, e := filepath.Glob("../plugins/*/manifest.json")
	if e != nil || len(files) != len(b) {
		t.Fatal("unregistered plugin directory", e)
	}
	for _, f := range files {
		raw, e := os.ReadFile(f)
		if e != nil {
			t.Fatal(e)
		}
		var m struct {
			ID string `json:"id"`
		}
		if e = json.Unmarshal(raw, &m); e != nil {
			t.Fatal(e)
		}
		if filepath.Base(filepath.Dir(f)) != m.ID {
			t.Fatal("directory ID drift")
		}
		if _, e = Get(m.ID); e != nil {
			t.Fatal(e)
		}
	}
}
