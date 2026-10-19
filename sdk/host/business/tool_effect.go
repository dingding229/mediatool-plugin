package business

import "strings"

// ToolEffect is shared by the SDK and host. Unknown/local scopes are classified
// conservatively; execution scopes never expand the instance's live grants.
func ToolEffect(scopes []string) string {
	rank := 0
	for _, s := range scopes {
		if s == "system.exec" || s == "system.shell" || strings.HasPrefix(s, "network") || strings.HasPrefix(s, "fs.media") || s == "auth.provider" || strings.HasPrefix(s, "events.chain.") || strings.HasPrefix(s, "module.") {
			rank = 2
		}
		if s == "plugin.data.write" || s == "plugin.config.write" || s == "notifications.send" {
			rank = max(rank, 1)
		}
		for _, m := range Methods() {
			if m.Scope != s {
				continue
			}
			if m.Level == "D" {
				rank = 2
			} else if m.Level == "W" {
				rank = max(rank, 1)
			}
		}
	}
	return []string{"read", "write", "dangerous"}[rank]
}
func ValidToolEffect(effect string, scopes []string) bool {
	need := ToolEffect(scopes)
	return effect == "dangerous" || effect == need || effect == "write" && need == "read"
}
