package contracts

import (
	"bytes"
	"encoding/json"
	"io"
	"unicode/utf8"
)

const MaxJSON = 4 << 20

// ValidateJSON rejects duplicate keys, invalid UTF-8, excessive depth and trailing JSON.
func ValidateJSON(b []byte) error {
	if len(b) == 0 || len(b) > MaxJSON || !utf8.Valid(b) {
		return NewError(InvalidParams)
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var walk func(int) error
	walk = func(depth int) error {
		if depth > 64 {
			return NewError(InvalidParams)
		}
		t, err := d.Token()
		if err != nil {
			return NewError(InvalidParams)
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return nil
		}
		switch delim {
		case '{':
			seen := map[string]bool{}
			for d.More() {
				k, e := d.Token()
				if e != nil {
					return NewError(InvalidParams)
				}
				key, ok := k.(string)
				if !ok || seen[key] {
					return NewError(InvalidParams)
				}
				seen[key] = true
				if e = walk(depth + 1); e != nil {
					return e
				}
			}
			end, e := d.Token()
			if e != nil || end != json.Delim('}') {
				return NewError(InvalidParams)
			}
		case '[':
			for d.More() {
				if e := walk(depth + 1); e != nil {
					return e
				}
			}
			end, e := d.Token()
			if e != nil || end != json.Delim(']') {
				return NewError(InvalidParams)
			}
		default:
			return NewError(InvalidParams)
		}
		return nil
	}
	if err := walk(0); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return NewError(InvalidParams)
	}
	return nil
}
func Decode(b []byte, out any) error {
	if err := ValidateJSON(b); err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err := d.Decode(out); err != nil {
		return NewError(InvalidParams)
	}
	return nil
}
