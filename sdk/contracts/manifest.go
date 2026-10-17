package contracts

import (
	"bytes"
	"encoding/json"
	"io"
	"path"
	"regexp"
	"strings"
)

type APIRoute struct {
	ID              string   `json:"id"`
	Path            string   `json:"path"`
	Methods         []string `json:"methods"`
	Auth            string   `json:"auth"`
	ExecutionScopes []string `json:"executionScopes"`
	Summary         string   `json:"summary,omitempty"`
}
type Permission struct {
	Scope     string              `json:"scope"`
	Resources map[string][]string `json:"resources,omitempty"`
}
type Extensions struct {
	Forms         bool       `json:"forms,omitempty"`
	Pages         []string   `json:"pages,omitempty"`
	Dashboards    []string   `json:"dashboards,omitempty"`
	API           []APIRoute `json:"api,omitempty"`
	Services      []string   `json:"services,omitempty"`
	Commands      []string   `json:"commands,omitempty"`
	Events        []string   `json:"events,omitempty"`
	Chains        []string   `json:"chains,omitempty"`
	Modules       []string   `json:"modules,omitempty"`
	Actions       []string   `json:"actions,omitempty"`
	AuthProviders []string   `json:"authProviders,omitempty"`
	AgentTools    []string   `json:"agentTools,omitempty"`
}
type Sidebar struct {
	Page  string `json:"page"`
	Title string `json:"title"`
	Icon  string `json:"icon,omitempty"`
}
type UI struct {
	Mode    string    `json:"mode"`
	Entry   string    `json:"entry,omitempty"`
	Sidebar []Sidebar `json:"sidebar,omitempty"`
}
type Resources struct {
	AddressSpaceMiB int `json:"addressSpaceMiB"`
	RSSMiB          int `json:"rssMiB"`
	CPUQuotaPercent int `json:"cpuQuotaPercent"`
	Processes       int `json:"processes"`
	OpenFiles       int `json:"openFiles"`
	DataMiB         int `json:"dataMiB"`
	TmpMiB          int `json:"tmpMiB"`
}
type Manifest struct {
	ContractVersion     int               `json:"contractVersion"`
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Version             string            `json:"version"`
	Author              string            `json:"author"`
	Description         string            `json:"description"`
	Level               int               `json:"level"`
	Icon                string            `json:"icon"`
	History             map[string]string `json:"history"`
	Permissions         []Permission      `json:"permissions"`
	Extensions          Extensions        `json:"extensions"`
	UI                  UI                `json:"ui"`
	Resources           Resources         `json:"resources"`
	ConfigSchemaVersion int               `json:"configSchemaVersion"`
	DataSchemaVersion   int               `json:"dataSchemaVersion"`
}

func ParseManifest(raw []byte) (Manifest, error) {
	var m Manifest
	m.Resources = Resources{8192, 256, 100, 64, 256, 256, 64}
	if len(raw) > 128<<10 || Decode(raw, &m) != nil {
		return m, NewError(InvalidParams)
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(raw, &fields) != nil {
		return m, NewError(InvalidParams)
	}
	for _, key := range []string{"contractVersion", "id", "name", "version", "author", "description", "level", "icon", "history", "permissions", "extensions", "ui", "resources", "configSchemaVersion", "dataSchemaVersion"} {
		if len(fields[key]) == 0 || string(fields[key]) == "null" {
			return m, NewError(InvalidParams)
		}
	}
	for _, key := range []string{"extensions", "ui", "resources", "history"} {
		if len(fields[key]) == 0 || fields[key][0] != '{' {
			return m, NewError(InvalidParams)
		}
	}
	tokens := json.NewDecoder(bytes.NewReader(raw))
	for {
		token, e := tokens.Token()
		if e == io.EOF {
			break
		}
		if e != nil || token == nil {
			return m, NewError(InvalidParams)
		}
	}
	validID, _ := regexp.MatchString(`^[a-z][a-z0-9-]{2,63}$`, m.ID)
	semver, _ := regexp.MatchString(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$`, m.Version)
	if !validID || !semver || m.ContractVersion != ContractVersion || m.Name == "" || len(m.Name) > 120 || m.Author == "" || len(m.Author) > 120 || len(m.Description) > 4000 || m.Level < 0 || m.Level > 9 || !safePath(m.Icon) || len(m.History) > 100 || m.History[m.Version] == "" || m.ConfigSchemaVersion < 1 || m.ConfigSchemaVersion > 10000 || m.DataSchemaVersion < 1 || m.DataSchemaVersion > 10000 || m.Permissions == nil || len(m.Permissions) > 128 {
		return m, NewError(InvalidParams)
	}
	seen := map[string]bool{}
	for _, p := range m.Permissions {
		if !KnownScope(p.Scope) || seen[p.Scope] {
			return m, NewError(InvalidParams)
		}
		seen[p.Scope] = true
		for k, ids := range p.Resources {
			switch k {
			case "libraryIds", "storageIds", "siteIds", "rootIds", "hosts", "eventTypes", "settingKeys", "workflowIds":
			default:
				return m, NewError(InvalidParams)
			}
			if !validIDs(ids) {
				return m, NewError(InvalidParams)
			}
		}
	}
	if m.UI.Mode != "json" && m.UI.Mode != "vue" {
		return m, NewError(InvalidParams)
	}
	if m.UI.Mode == "vue" && (!strings.HasPrefix(m.UI.Entry, "ui/dist/") || !safePath(m.UI.Entry) || (!strings.HasSuffix(m.UI.Entry, ".js") && !strings.HasSuffix(m.UI.Entry, ".mjs")) || len(m.UI.Entry) > 200 || !seen["ui.trusted"]) {
		return m, NewError(InvalidParams)
	}
	if len(m.UI.Sidebar) > 8 || (!seen["ui.sidebar"] && len(m.UI.Sidebar) > 0) {
		return m, NewError(InvalidParams)
	}
	if (m.Extensions.Forms || len(m.Extensions.Pages) > 0 || len(m.Extensions.Dashboards) > 0) && m.UI.Mode == "json" && !seen["ui.json"] {
		return m, NewError(InvalidParams)
	}
	for _, binding := range []struct {
		n     int
		scope string
	}{{len(m.Extensions.API), "api.register"}, {len(m.Extensions.Services), "services.register"}, {len(m.Extensions.Commands), "commands.register"}, {len(m.Extensions.Actions), "actions.register"}, {len(m.Extensions.AuthProviders), "auth.provider"}, {len(m.Extensions.AgentTools), "agent.tools"}} {
		if binding.n > 0 && !seen[binding.scope] {
			return m, NewError(InvalidParams)
		}
	}
	for _, binding := range []struct {
		ids    []string
		prefix string
	}{{m.Extensions.Events, "events.subscribe."}, {m.Extensions.Chains, "events.chain."}, {m.Extensions.Modules, "module."}} {
		for _, id := range binding.ids {
			if !KnownScope(binding.prefix+id) || !seen[binding.prefix+id] {
				return m, NewError(InvalidParams)
			}
		}
	}
	for _, side := range m.UI.Sidebar {
		if side.Page == "" || len(side.Page) > 64 || side.Title == "" || len(side.Title) > 120 || len(side.Icon) > 64 {
			return m, NewError(InvalidParams)
		}
		found := false
		for _, id := range m.Extensions.Pages {
			if id == side.Page {
				found = true
			}
		}
		if !found {
			return m, NewError(InvalidParams)
		}
	}
	for _, ids := range [][]string{m.Extensions.Pages, m.Extensions.Dashboards, m.Extensions.Services, m.Extensions.Commands, m.Extensions.Events, m.Extensions.Chains, m.Extensions.Modules, m.Extensions.Actions, m.Extensions.AuthProviders, m.Extensions.AgentTools} {
		if !validIDs(ids) {
			return m, NewError(InvalidParams)
		}
	}
	for _, route := range m.Extensions.API {
		if route.ID == "" || len(route.ID) > 64 || route.Path == "" || len(route.Path) > 200 || len(route.Summary) > 500 || len(route.Methods) < 1 || len(route.Methods) > 5 || !validIDs(route.Methods) || !validIDs(route.ExecutionScopes) || route.ExecutionScopes == nil {
			return m, NewError(InvalidParams)
		}
		for _, scope := range route.ExecutionScopes {
			if !seen[scope] {
				return m, NewError(InvalidParams)
			}
		}
		if route.Auth == "public" && !seen["api.public"] {
			return m, NewError(InvalidParams)
		}
		switch route.Auth {
		case "admin", "apikey", "bearer", "public":
		default:
			return m, NewError(InvalidParams)
		}
		for _, method := range route.Methods {
			if route.Auth == "public" && method != "GET" && !seen["api.public.write"] {
				return m, NewError(InvalidParams)
			}
			switch method {
			case "GET", "POST", "PUT", "PATCH", "DELETE":
			default:
				return m, NewError(InvalidParams)
			}
		}
	}
	if len(m.Extensions.API) > 32 {
		return m, NewError(InvalidParams)
	}
	for version, note := range m.History {
		if len(version) > 64 || len(note) > 4000 {
			return m, NewError(InvalidParams)
		}
	}
	r := m.Resources
	if r.AddressSpaceMiB < 256 || r.AddressSpaceMiB > 8192 || r.RSSMiB < 64 || r.RSSMiB > 1024 || r.CPUQuotaPercent < 10 || r.CPUQuotaPercent > 200 || r.Processes < 1 || r.Processes > 128 || r.OpenFiles < 32 || r.OpenFiles > 1024 || r.DataMiB < 1 || r.DataMiB > 2048 || r.TmpMiB < 1 || r.TmpMiB > 256 {
		return m, NewError(InvalidParams)
	}
	return m, nil
}
func safePath(s string) bool {
	valid, _ := regexp.MatchString(`^[A-Za-z0-9_][A-Za-z0-9_./-]*$`, s)
	return valid && s != "" && len(s) <= 240 && path.Clean(s) == s && !strings.ContainsAny(s, "\\\x00:") && !strings.Contains(s, "..") && !strings.HasPrefix(s, "/")
}
func validIDs(ids []string) bool {
	if len(ids) > 64 {
		return false
	}
	seen := map[string]bool{}
	for _, id := range ids {
		if id == "" || len(id) > 200 || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}
