package business_test

import (
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	b "github.com/dingding229/mediatool-plugin/sdk/host/business"
	"github.com/dingding229/mediatool-plugin/sdk/host/library"
	"reflect"
	"testing"
)

func TestI12ArchiveCapabilitiesRemoved(t *testing.T) {
	for _, name := range []string{"library.archive.list", "library.archive.preview", "library.archive.confirmCleanup"} {
		if _, ok := b.Lookup(name); ok {
			t.Fatalf("removed method remains: %s", name)
		}
	}
	for _, scope := range []string{"library.archive.read", "library.archive.delete"} {
		if contracts.KnownScope(scope) {
			t.Fatalf("removed scope remains: %s", scope)
		}
	}
	client := reflect.TypeOf((*library.Client)(nil))
	for _, name := range []string{"Archives", "ArchivePreview", "ArchiveCleanup"} {
		if _, ok := client.MethodByName(name); ok {
			t.Fatalf("removed typed method remains: %s", name)
		}
	}
}
