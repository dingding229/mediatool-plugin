package contracts

import (
	"encoding/json"
	"testing"
)

func TestToolSchemaBoundsAndRejectsUnknownKeywords(t *testing.T) {
	s := json.RawMessage(`{"type":"object","required":["count"],"properties":{"count":{"type":"integer","minimum":1,"maximum":3}},"additionalProperties":false}`)
	for _, raw := range []string{`{"count":1}`, `{"count":3}`} {
		if e := ValidateToolInput(s, []byte(raw)); e != nil {
			t.Fatal(raw, e)
		}
	}
	for _, raw := range []string{`{}`, `{"count":0}`, `{"count":1.2}`, `{"count":1,"extra":true}`, `{"count":1,"count":2}`} {
		if ValidateToolInput(s, []byte(raw)) == nil {
			t.Fatal("invalid accepted", raw)
		}
	}
	for _, raw := range []string{`{"$ref":"https://remote.invalid/schema"}`, `{"type":"string","pattern":".*"}`, `{"type":"array"}`, `{"type":"string","minLength":3,"maxLength":1}`} {
		if _, e := ParseToolSchema([]byte(raw)); e == nil {
			t.Fatal("unsupported schema accepted", raw)
		}
	}
}
