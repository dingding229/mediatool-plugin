// Package schema owns the sole manifest schema bundled with the pinned SDK.
package schema

import _ "embed"

//go:embed manifest-v2.schema.json
var manifest string

const Fingerprint = "3629df8254b7d46d9319cc1301f9f0961c9207c45c426c381c07e8b66b86fc4a"

func JSON() []byte { return []byte(manifest) }
