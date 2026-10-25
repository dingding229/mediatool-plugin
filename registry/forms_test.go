package registry_test

import (
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	hello "github.com/dingding229/mediatool-plugin/plugins/hello-go"
	refresh "github.com/dingding229/mediatool-plugin/plugins/library-refresh-notify"
	cleanup "github.com/dingding229/mediatool-plugin/plugins/seeding-cleanup"
	checkin "github.com/dingding229/mediatool-plugin/plugins/site-checkin"
	assistant "github.com/dingding229/mediatool-plugin/plugins/torrent-assistant"
	watch "github.com/dingding229/mediatool-plugin/plugins/watch-organize"
	wishlist "github.com/dingding229/mediatool-plugin/plugins/wishlist-subscribe"
	"github.com/dingding229/mediatool-plugin/registry"
	"github.com/dingding229/mediatool-plugin/sdk"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"github.com/dingding229/mediatool-plugin/sdk/events"
	"github.com/dingding229/mediatool-plugin/sdk/plugintest"
)

var configSchemas = map[string]any{"hello-go": hello.Config{}, "library-refresh-notify": refresh.Config{}, "seeding-cleanup": cleanup.Config{}, "site-checkin": checkin.Config{}, "torrent-assistant": assistant.Config{}, "watch-organize": watch.Config{}, "wishlist-subscribe": wishlist.Config{}}

func controls(nodes []contracts.Component) map[string]contracts.Component {
	out := map[string]contracts.Component{}
	for _, node := range nodes {
		if node.Model != "" {
			out[node.Model] = node
		}
		for k, v := range controls(node.Children) {
			out[k] = v
		}
	}
	return out
}
func number(t *testing.T, n contracts.Component, key string) float64 {
	t.Helper()
	var value float64
	if json.Unmarshal(n.Props[key], &value) != nil {
		t.Fatalf("missing %s on %s", key, n.Model)
	}
	return value
}
func TestEveryBuiltinConfigSchemaFieldIsEditableWithoutLifecycleOrHost(t *testing.T) {
	defs, err := registry.All()
	if err != nil || len(defs) != len(configSchemas) {
		t.Fatal(err, "unregistered configuration schema")
	}
	for _, def := range defs {
		t.Run(def.Manifest.ID, func(t *testing.T) {
			// No InitPlugin, Host, grants, process or persisted configuration: same pure
			// build-time path used to export stopped-instance schemas.
			p := def.Factory(sdk.Services{})
			form, err := p.(sdk.FormProvider).GetForm(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			fields := controls(form.Content)
			schema := reflect.TypeOf(configSchemas[def.Manifest.ID])
			if schema == nil {
				t.Fatal("missing config schema")
			}
			for i := 0; i < schema.NumField(); i++ {
				field := schema.Field(i)
				key := strings.Split(field.Tag.Get("json"), ",")[0]
				n, ok := fields[key]
				if !ok {
					t.Fatal("configuration cannot be set in form", key)
				}
				switch field.Type.Kind() {
				case reflect.Bool:
					if n.Component != "VSwitch" && n.Component != "VCheckbox" {
						t.Fatal(key, "boolean control missing")
					}
				case reflect.String:
					if n.Component != "VTextField" && n.Component != "VTextarea" && n.Component != "VSelect" && n.Component != "VCronInput" {
						t.Fatal(key, "text control missing")
					}
				default:
					if n.Component != "VNumberInput" {
						t.Fatal(key, "number must not be narrowed to enum/text")
					}
					if number(t, n, "min") > number(t, n, "max") || number(t, n, "step") <= 0 {
						t.Fatal(key, "invalid interval")
					}
				}
			}
		})
	}
	// Map-backed example configuration must also declare every executable key in
	// its finite schema; adding a handler flag cannot silently escape the guard.
	files, _ := filepath.Glob("../plugins/hello-go/*.go")
	schema := reflect.TypeOf(hello.Config{})
	known := map[string]bool{}
	for i := 0; i < schema.NumField(); i++ {
		known[schema.Field(i).Tag.Get("json")] = true
	}
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		tree, e := parser.ParseFile(token.NewFileSet(), file, nil, 0)
		if e != nil {
			t.Fatal(e)
		}
		ast.Inspect(tree, func(n ast.Node) bool {
			x, ok := n.(*ast.IndexExpr)
			if !ok {
				return true
			}
			sel, ok := x.X.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "config" {
				return true
			}
			literal, ok := x.Index.(*ast.BasicLit)
			if ok && literal.Kind == token.STRING {
				key, _ := strconv.Unquote(literal.Value)
				if !known[key] {
					t.Error("executable configuration absent from schema", key)
				}
			}
			return true
		})
	}
}
func TestNumericFormBoundariesMatchAcceptedConfiguration(t *testing.T) {
	// Last value is 1 for integer-only settings, 0 for continuous numeric ranges.
	bounds := map[string]map[string][3]float64{
		"library-refresh-notify": {"intervalSeconds": {60, 3600, 1}},
		"seeding-cleanup":        {"minRatio": {0, 10000, 0}, "minSeedingHours": {1, 87600, 0}, "batchSize": {1, 20, 1}},
		"torrent-assistant":      {"downloadKiB": {0, 102400, 0}, "uploadKiB": {0, 102400, 0}, "ratio": {0, 1000, 0}, "seedHours": {0, 87600, 0}, "batch": {1, 10, 1}, "stallMinutes": {5, 1440, 0}},
		"watch-organize":         {"depth": {1, 5, 1}, "quietSeconds": {30, 3600, 1}},
		"wishlist-subscribe":     {"intervalMinutes": {30, 1440, 1}, "perRun": {1, 20, 1}},
	}
	for id, values := range bounds {
		def, _ := registry.Get(id)
		form, err := def.Factory(sdk.Services{}).(sdk.FormProvider).GetForm(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		fields := controls(form.Content)
		for key, bound := range values {
			t.Run(id+"/"+key, func(t *testing.T) {
				n := fields[key]
				if number(t, n, "min") != bound[0] || number(t, n, "max") != bound[1] {
					t.Fatal("form narrows legal interval", n)
				}
				var strict bool
				json.Unmarshal(n.Props["stepStrictly"], &strict)
				if strict != (bound[2] == 1) {
					t.Fatal("step snapping would truncate legal fractions")
				}
				middle := bound[0] + .001
				if strict {
					middle = bound[0] + 1
				}
				for _, sample := range []struct {
					value float64
					valid bool
				}{{bound[0], true}, {middle, true}, {bound[1], true}, {bound[0] - 1, false}, {bound[1] + 1, false}} {
					var config map[string]any
					json.Unmarshal(form.Defaults, &config)
					config[key] = sample.value
					raw, _ := json.Marshal(config)
					scopes := []string{}
					for _, permission := range def.Manifest.Permissions {
						scopes = append(scopes, permission.Scope)
					}
					host := plugintest.NewHost(scopes)
					p := def.Factory(sdk.Services{Host: host.Client(), Events: events.New()})
					err := p.InitPlugin(context.Background(), sdk.Config{Value: raw})
					if (err == nil) != sample.valid {
						t.Fatalf("value %v valid=%v: %v", sample.value, sample.valid, err)
					}
					if err == nil {
						p.StopService(context.Background())
					}
				}
			})
		}
	}
}

func TestFormTitlesAndFiniteEnumsArePresentationOnly(t *testing.T) {
	defs, e := registry.All()
	if e != nil {
		t.Fatal(e)
	}
	want := map[string][]string{"library-refresh-notify.provider": {"", "builtin", "emby"}, "site-checkin.method": {"", "GET", "POST"}, "wishlist-subscribe.source": {"tmdb-list", "douban-wish", "douban-subject"}}
	seen := 0
	for _, def := range defs {
		form, e := def.Factory(sdk.Services{}).(sdk.FormProvider).GetForm(context.Background())
		if e != nil {
			t.Fatal(e)
		}
		for key, n := range controls(form.Content) {
			var title, label string
			if json.Unmarshal(n.Props["title"], &title) != nil || title == "" || json.Unmarshal(n.Props["label"], &label) != nil || title != label {
				t.Fatal(def.Manifest.ID, key, "missing formal title")
			}
			if values, ok := want[def.Manifest.ID+"."+key]; ok {
				var items, enum []string
				if json.Unmarshal(n.Props["items"], &items) != nil || json.Unmarshal(n.Props["enum"], &enum) != nil || !reflect.DeepEqual(items, values) || !reflect.DeepEqual(enum, values) || n.Component != "VSelect" && n.Component != "VCronInput" {
					t.Fatal(def.Manifest.ID, key, "enum contract")
				}
				seen++
			}
		}
	}
	if seen != len(want) {
		t.Fatal("missing enum fields", seen)
	}
}
