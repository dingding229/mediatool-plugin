// Package registry lists the built-in library. It never initializes an instance.
package registry

import (
	"github.com/dingding229/mediatool-plugin/plugins/hello-go"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
)

type Definition struct {
	Manifest contracts.Manifest
	Factory  sdk.Factory
}

// No init functions, IO, or plugin factories run while importing this package.
func All() ([]Definition, error) {
	raw := []struct {
		id       string
		metadata []byte
		factory  sdk.Factory
	}{{"hello-go", hello.ManifestJSON(), hello.New}}
	out := make([]Definition, 0, len(raw))
	seen := map[string]bool{}
	for _, entry := range raw {
		m, e := contracts.ParseManifest(entry.metadata)
		if e != nil {
			return nil, e
		}
		if m.ID != entry.id || seen[m.ID] || entry.factory == nil {
			return nil, contracts.NewError(contracts.InvalidParams)
		}
		seen[m.ID] = true
		out = append(out, Definition{m, entry.factory})
	}
	return out, nil
}
func Get(id string) (Definition, error) {
	defs, e := All()
	if e != nil {
		return Definition{}, e
	}
	for _, d := range defs {
		if d.Manifest.ID == id {
			return d, nil
		}
	}
	return Definition{}, contracts.NewError(contracts.NotFound)
}
