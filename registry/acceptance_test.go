package registry_test

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/registry"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"testing"
)

func TestAcceptanceDescriptorsGoldenAndNoInitEffects(t *testing.T) {
	for _, name := range []string{"site-checkin", "library-refresh-notify", "seeding-cleanup", "watch-organize", "torrent-assistant", "wishlist-subscribe"} {
		t.Run(name, func(t *testing.T) {
			d, e := registry.Get(name)
			if e != nil {
				t.Fatal(e)
			}
			scopes := []string{}
			for _, p := range d.Manifest.Permissions {
				scopes = append(scopes, p.Scope)
			}
			mock := plugintest.NewHost(scopes)
			identity := sdk.Identity{InstanceID: "fixture", DefinitionFingerprint: "definition", CatalogBuild: "build"}
			h, e := plugintest.Start(context.Background(), d.Manifest, d.Factory, mock, identity)
			if e != nil {
				t.Fatal(e)
			}
			defer h.Close()
			q := contracts.Initialize{ContractVersion: 2, InstanceID: "fixture", Generation: "generation", ConfigRevision: 1, DefinitionFingerprint: "definition", CatalogBuild: "build", Capabilities: scopes}
			if _, e = h.Initialize(context.Background(), q); e != nil {
				t.Fatal(e)
			}
			if e = h.Init(context.Background(), contracts.Config{Value: json.RawMessage(`{}`), Revision: 1}); e != nil {
				t.Fatal(e)
			}
			var snap contracts.Snapshot
			if e = h.Call(context.Background(), "plugin.describe", nil, &snap); e != nil {
				t.Fatal(e)
			}
			if e = contracts.ValidateSnapshot(d.Manifest, scopes, snap); e != nil {
				t.Fatal(e)
			}
			if len(mock.Calls()) != 0 {
				t.Fatal("inactive instance executed Host calls", mock.Calls())
			}
			plugintest.Golden(t, "../plugins/"+name+"/testdata/manifest.json", d.Manifest)
			plugintest.Golden(t, "../plugins/"+name+"/testdata/descriptor.json", snap)
		})
	}
}
