package plugintest

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

// Golden never rewrites fixtures automatically; review expected contract changes explicitly.
func Golden(t testing.TB, path string, value any) {
	t.Helper()
	b, e := json.MarshalIndent(value, "", "  ")
	if e != nil {
		t.Fatal(e)
	}
	b = append(b, '\n')
	want, e := os.ReadFile(path)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(b, want) {
		t.Errorf("golden mismatch %s\nactual:\n%s", path, b)
	}
}
