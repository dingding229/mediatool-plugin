package contracts

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"regexp"
	"unicode/utf8"
)

const ToolInputLimit = 16 << 10
const ToolOutputLimit = 32 << 10

// ToolSchema is a bounded JSON Schema subset. Unknown keywords/$ref are rejected,
// never silently ignored or remotely resolved. No coercion or default injection.
type ToolSchema struct {
	Type                 string                `json:"type"`
	Title                string                `json:"title,omitempty"`
	Description          string                `json:"description,omitempty"`
	Properties           map[string]ToolSchema `json:"properties,omitempty"`
	Required             []string              `json:"required,omitempty"`
	AdditionalProperties *bool                 `json:"additionalProperties,omitempty"`
	Items                *ToolSchema           `json:"items,omitempty"`
	Enum                 []json.RawMessage     `json:"enum,omitempty"`
	Minimum              *float64              `json:"minimum,omitempty"`
	Maximum              *float64              `json:"maximum,omitempty"`
	MinLength            *int                  `json:"minLength,omitempty"`
	MaxLength            *int                  `json:"maxLength,omitempty"`
	MinItems             *int                  `json:"minItems,omitempty"`
	MaxItems             *int                  `json:"maxItems,omitempty"`
}

func ParseToolSchema(raw json.RawMessage) (ToolSchema, error) {
	var s ToolSchema
	if len(raw) > 8<<10 || Decode(raw, &s) != nil {
		return s, NewError(InvalidParams)
	}
	count := 0
	var walk func(ToolSchema, int) bool
	walk = func(x ToolSchema, depth int) bool {
		count++
		if depth > 8 || count > 128 || len(x.Description) > 2000 || len(x.Title) > 120 || len(x.Enum) > 32 {
			return false
		}
		switch x.Type {
		case "object":
		case "array":
			if x.Items == nil {
				return false
			}
		case "string", "number", "integer", "boolean", "null":
		default:
			return false
		}
		if len(x.Properties) > 64 || len(x.Required) > 64 {
			return false
		}
		seen := map[string]bool{}
		for _, k := range x.Required {
			if _, ok := x.Properties[k]; !ok || seen[k] {
				return false
			}
			seen[k] = true
		}
		for k, p := range x.Properties {
			if len(k) > 100 || !walk(p, depth+1) {
				return false
			}
		}
		if x.Items != nil && !walk(*x.Items, depth+1) {
			return false
		}
		for _, v := range x.Enum {
			if ValidateJSON(v) != nil {
				return false
			}
		}
		for _, pair := range [][2]*int{{x.MinLength, x.MaxLength}, {x.MinItems, x.MaxItems}} {
			if pair[0] != nil && (*pair[0] < 0 || *pair[0] > 16384) || pair[1] != nil && (*pair[1] < 0 || *pair[1] > 16384) || pair[0] != nil && pair[1] != nil && *pair[0] > *pair[1] {
				return false
			}
		}
		return !(x.Minimum != nil && x.Maximum != nil && *x.Minimum > *x.Maximum)
	}
	if !walk(s, 0) {
		return s, NewError(InvalidParams)
	}
	return s, nil
}
func ValidateToolInput(schema, input json.RawMessage) error {
	s, e := ParseToolSchema(schema)
	if e != nil {
		return e
	}
	if len(input) > ToolInputLimit || ValidateJSON(input) != nil {
		return NewError(InvalidParams)
	}
	var value any
	d := json.NewDecoder(bytes.NewReader(input))
	d.UseNumber()
	if d.Decode(&value) != nil {
		return NewError(InvalidParams)
	}
	var check func(ToolSchema, any) bool
	bounds := func(n int, lo, hi *int) bool { return (lo == nil || n >= *lo) && (hi == nil || n <= *hi) }
	check = func(x ToolSchema, v any) bool {
		if x.Enum != nil {
			found := false
			for _, b := range x.Enum {
				var other any
				d := json.NewDecoder(bytes.NewReader(b))
				d.UseNumber()
				if d.Decode(&other) == nil && reflect.DeepEqual(v, other) {
					found = true
				}
			}
			if !found {
				return false
			}
		}
		switch x.Type {
		case "object":
			m, ok := v.(map[string]any)
			if !ok {
				return false
			}
			for _, k := range x.Required {
				if _, ok := m[k]; !ok {
					return false
				}
			}
			for k, v := range m {
				p, ok := x.Properties[k]
				if !ok {
					if x.AdditionalProperties != nil && !*x.AdditionalProperties {
						return false
					}
					continue
				}
				if !check(p, v) {
					return false
				}
			}
		case "array":
			a, ok := v.([]any)
			if !ok || !bounds(len(a), x.MinItems, x.MaxItems) {
				return false
			}
			for _, v := range a {
				if !check(*x.Items, v) {
					return false
				}
			}
		case "string":
			s, ok := v.(string)
			if !ok || !bounds(utf8.RuneCountInString(s), x.MinLength, x.MaxLength) {
				return false
			}
		case "integer", "number":
			n, ok := v.(json.Number)
			if !ok {
				return false
			}
			f, e := n.Float64()
			if e != nil || math.IsInf(f, 0) || math.IsNaN(f) || x.Type == "integer" && math.Trunc(f) != f || x.Minimum != nil && f < *x.Minimum || x.Maximum != nil && f > *x.Maximum {
				return false
			}
		case "boolean":
			if _, ok := v.(bool); !ok {
				return false
			}
		case "null":
			if v != nil {
				return false
			}
		}
		return true
	}
	if !check(s, value) {
		return NewError(InvalidParams)
	}
	return nil
}
func ValidTool(d Descriptor) bool {
	ok, _ := regexp.MatchString(`^[a-z][a-z0-9_.-]{0,63}$`, d.ID)
	if !ok || d.Name != d.ID || d.Description == "" || len(d.RequiredScopes) > 32 {
		return false
	}
	if d.EffectKind != "read" && d.EffectKind != "write" && d.EffectKind != "dangerous" {
		return false
	}
	_, e := ParseToolSchema(d.InputSchema)
	return e == nil
}
