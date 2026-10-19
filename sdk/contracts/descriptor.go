package contracts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"
)

// ValidateSnapshot is shared by SDK and host; Host must supply live granted scopes.
func ValidateSnapshot(m Manifest, granted []string, s Snapshot) error {
	if s.Version != DescriptorVersion || len(s.Descriptors) > 128 {
		return NewError(UnsupportedVersion)
	}
	ids := map[string]bool{}
	paths := map[string]bool{}
	commands := map[string]bool{}
	for _, d := range s.Descriptors {
		if d.Handler == "" || len(d.Handler) > 128 || len(d.ID) > 200 || d.ID == "" || ids[d.Kind+":"+d.ID] || d.Handler != HandlerID(d.Kind, d.ID) || len(d.Name) > 120 || len(d.Summary) > 500 || len(d.Description) > 2000 || len(d.Category) > 120 || d.Priority < -1000 || d.Priority > 1000 {
			return NewError(InvalidExtensionResult)
		}
		ids[d.Kind+":"+d.ID] = true
		if d.ReplayWindowSeconds < 0 || d.ReplayWindowSeconds > 7*86400 || d.Kind != "event" && d.ReplayWindowSeconds != 0 {
			return NewError(InvalidParams)
		}
		required := ""
		allowed := []string{}
		switch d.Kind {
		case "api":
			required = "api.register"
			found := false
			for _, a := range m.Extensions.API {
				if a.ID == d.ID && a.Path == d.Path && a.Auth == d.Auth && slices.Equal(a.Methods, d.Methods) {
					found = true
					for _, scope := range a.ExecutionScopes {
						if !slices.Contains(d.RequiredScopes, scope) {
							return NewError(PermissionDenied)
						}
					}
				}
			}
			if !found || !ValidAPIPath(d.Path) {
				return NewError(InvalidExtensionResult)
			}
			for _, method := range d.Methods {
				key := method + ":" + d.Path
				if paths[key] {
					return NewError(Conflict)
				}
				paths[key] = true
				if d.Auth == "public" && (!slices.Contains(d.RequiredScopes, "api.public") || (method != "GET" && !slices.Contains(d.RequiredScopes, "api.public.write"))) {
					return NewError(PermissionDenied)
				}
			}
		case "service":
			required = "services.register"
			allowed = m.Extensions.Services
			switch d.Trigger {
			case "cron":
				if len(d.Cron) > 120 || len(strings.Fields(d.Cron)) != 5 {
					return NewError(InvalidParams)
				}
			case "interval":
				if d.IntervalMS < 1000 || d.IntervalMS > 86400000*365 {
					return NewError(InvalidParams)
				}
			case "date":
				if _, e := time.Parse(time.RFC3339Nano, d.At); e != nil {
					return NewError(InvalidParams)
				}
			case "once":
			default:
				return NewError(InvalidParams)
			}
			if len(d.ScheduleSetting) > 64 {
				return NewError(InvalidParams)
			}
		case "command":
			required = "commands.register"
			allowed = m.Extensions.Commands
			ok, _ := regexp.MatchString(`^[a-z][a-z0-9_-]{0,63}$`, d.Cmd)
			if !ok || d.Cmd != d.ID || commands[d.Cmd] {
				return NewError(Conflict)
			}
			commands[d.Cmd] = true
		case "tool":
			required = "agent.tools"
			allowed = m.Extensions.AgentTools
			if !ValidTool(d) {
				return NewError(InvalidExtensionResult)
			}
		case "action":
			required = "actions.register"
			allowed = m.Extensions.Actions
		case "module":
			required = "module." + d.ID
			allowed = m.Extensions.Modules
		case "chain":
			required = "events.chain." + d.ID
			allowed = m.Extensions.Chains
			if d.ID == "auth.verification" && (!slices.Contains(m.Extensions.AuthProviders, d.ProviderID) || !slices.Contains(d.RequiredScopes, "auth.provider") || d.Issuer == "" || len(d.Issuer) > 200 || d.Audience == "" || len(d.Audience) > 200 || d.Origin == "" || len(d.Origin) > 200) {
				return NewError(InvalidExtensionResult)
			}
		case "ui":
			required = "ui.json"
			if m.Extensions.Forms {
				allowed = append(allowed, "form")
			}
			if len(m.Extensions.Dashboards) > 0 {
				allowed = append(allowed, "meta")
			}
			for _, k := range m.Extensions.Pages {
				allowed = append(allowed, "page."+k)
			}
			for _, k := range m.Extensions.Dashboards {
				allowed = append(allowed, "dashboard."+k)
			}
		case "event":
			required = "events.subscribe." + d.ID
			allowed = m.Extensions.Events
		default:
			return NewError(InvalidExtensionResult)
		}
		if d.Kind != "api" && !slices.Contains(allowed, d.ID) {
			return NewError(PermissionDenied)
		}
		if !slices.Contains(d.RequiredScopes, required) {
			return NewError(PermissionDenied)
		}
		for _, scope := range d.RequiredScopes {
			declared := false
			for _, p := range m.Permissions {
				if p.Scope == scope {
					declared = true
					break
				}
			}
			if !declared || !slices.Contains(granted, scope) {
				return NewError(PermissionDenied)
			}
		}
		for _, raw := range []json.RawMessage{d.Kwargs, d.Data} {
			if len(raw) > 16<<10 || (len(raw) > 0 && ValidateJSON(raw) != nil) {
				return NewError(InvalidParams)
			}
		}
	}
	return nil
}
func ValidAPIPath(p string) bool {
	if len(p) > 200 || !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || path.Clean(p) != p || strings.ContainsAny(p, "?%#\\{}*\x00\r\n ") {
		return false
	}
	return true
}

func HandlerID(kind, id string) string {
	b := sha256.Sum256([]byte(kind + "\x00" + id))
	return kind + "." + hex.EncodeToString(b[:12])
}
