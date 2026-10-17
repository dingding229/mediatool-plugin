package schema

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestSchemaFingerprint(t *testing.T) {
	if fmt.Sprintf("%x", sha256.Sum256(JSON())) != Fingerprint {
		t.Fatal("contract fingerprint drift")
	}
	b := JSON()
	b[0] = '!'
	if JSON()[0] != '{' {
		t.Fatal("schema accessor shares mutable storage")
	}
}
