package registry_test

import (
	"context"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/registry"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
	"os"
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
			// Manifest declarations must have real bindings, not merely a golden
			// that accidentally blesses an omitted optional Go interface.
			found := map[string]bool{}
			for _, x := range snap.Descriptors {
				found[x.Kind+":"+x.ID] = true
			}
			for kind, ids := range map[string][]string{"action": d.Manifest.Extensions.Actions, "service": d.Manifest.Extensions.Services, "command": d.Manifest.Extensions.Commands} {
				for _, id := range ids {
					if !found[kind+":"+id] {
						t.Fatal("missing declared binding", kind, id)
					}
				}
			}
			if len(mock.Calls()) != 0 {
				t.Fatal("inactive instance executed Host calls", mock.Calls())
			}
			plugintest.Golden(t, "../plugins/"+name+"/testdata/manifest.json", d.Manifest)
			plugintest.Golden(t, "../plugins/"+name+"/testdata/descriptor.json", snap)
		})
	}
}

func TestManifestOperationTitlesReachRuntimeDescriptors(t *testing.T) {
	d, e := registry.Get("torrent-assistant")
	if e != nil {
		t.Fatal(e)
	}
	if d.Manifest.Extensions.Titles["command:downloader-status"] != "下载器状态" || d.Manifest.Extensions.Titles["service:assist"] != "下载辅助检查" {
		t.Fatal(d.Manifest.Extensions.Titles)
	}
	// The reviewed descriptor golden above exercises actual SDK registration,
	// including command providers that do not have their own Name member.
	raw, e := os.ReadFile("../plugins/torrent-assistant/testdata/descriptor.json")
	if e != nil {
		t.Fatal(e)
	}
	var snap contracts.Snapshot
	if e = json.Unmarshal(raw, &snap); e != nil {
		t.Fatal(e)
	}
	for key, title := range d.Manifest.Extensions.Titles {
		found := false
		for _, v := range snap.Descriptors {
			if v.Kind+":"+v.ID == key {
				found = true
				if v.Name != title {
					t.Fatal(key, v.Name, title)
				}
			}
		}
		if !found {
			t.Fatal(key)
		}
	}
}
