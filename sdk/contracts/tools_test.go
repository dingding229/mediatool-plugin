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
func TestToolNumbersDoNotRoundFractionalIntegersOrEnumerations(t *testing.T) {
	if ValidateToolInput([]byte(`{"type":"integer"}`), []byte(`9007199254740992.1`)) == nil {
		t.Fatal("fraction rounded into integer")
	}
	if ValidateToolInput([]byte(`{"type":"number","enum":[1]}`), []byte(`1.0`)) != nil {
		t.Fatal("numeric enum not semantic")
	}
	if ValidateToolInput([]byte(`{"type":"number"}`), []byte(`1e-999999999`)) == nil {
		t.Fatal("unbounded exponent")
	}
}
