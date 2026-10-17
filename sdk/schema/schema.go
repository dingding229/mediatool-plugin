// Package schema owns the sole manifest schema bundled with the pinned SDK.
package schema

import _ "embed"

//go:embed manifest-v2.schema.json
var manifest string

const Fingerprint = "a5fb191f1fef39c569218c11fba3b8c8544aacd1ded29cb7bd2ed9224003c97f"

func JSON() []byte { return []byte(manifest) }
