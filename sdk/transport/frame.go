// Package transport implements bounded, bidirectional JSON-RPC 2.0 over framed IO.
package transport

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/dingding229/mediatool-plugin/sdk/contracts"
	"io"
	"strconv"
	"strings"
)

type message struct {
	Version string           `json:"jsonrpc"`
	ID      string           `json:"id,omitempty"`
	Method  string           `json:"method,omitempty"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *contracts.Error `json:"error,omitempty"`
}

func validID(id, prefix string) bool {
	if !strings.HasPrefix(id, prefix+":") {
		return false
	}
	n, e := strconv.ParseUint(strings.TrimPrefix(id, prefix+":"), 10, 64)
	return e == nil && n > 0 && id == prefix+":"+strconv.FormatUint(n, 10)
}
func readFrame(r io.Reader, max int) (message, error) {
	var h [4]byte
	if _, err := io.ReadFull(r, h[:]); err != nil {
		return message{}, err
	}
	n := binary.BigEndian.Uint32(h[:])
	if n == 0 || uint64(n) > uint64(max) {
		return message{}, contracts.NewError(contracts.PayloadTooLarge)
	}
	b := make([]byte, int(n))
	if _, err := io.ReadFull(r, b); err != nil {
		return message{}, err
	}
	if trimmed := bytes.TrimSpace(b); len(trimmed) == 0 || trimmed[0] != '{' {
		return message{}, contracts.NewError(contracts.InvalidRequest)
	}
	var m message
	if err := contracts.Decode(b, &m); err != nil {
		return m, err
	}
	var fields map[string]json.RawMessage
	if json.Unmarshal(b, &fields) != nil {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	if id, exists := fields["id"]; exists && (m.ID == "" || bytes.Equal(bytes.TrimSpace(id), []byte("null"))) {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	if _, exists := fields["method"]; exists && m.Method == "" {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	_, hasResult := fields["result"]
	_, hasError := fields["error"]
	if m.Method != "" && (hasResult || hasError) {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	if m.Method == "" && hasResult == hasError {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	if m.Version != "2.0" || len(m.Method) > 128 {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	if m.Method != "" {
		if m.Result != nil || m.Error != nil {
			return m, contracts.NewError(contracts.InvalidRequest)
		}
	} else if m.ID == "" || m.Params != nil || (m.Result == nil) == (m.Error == nil) {
		return m, contracts.NewError(contracts.InvalidRequest)
	}
	return m, nil
}
func encode(m message, max int) ([]byte, error) {
	b, e := safeMarshal(m)
	if e != nil {
		return nil, contracts.NewError(contracts.InternalError)
	}
	if contracts.ValidateJSON(b) != nil {
		return nil, contracts.NewError(contracts.InvalidParams)
	}
	if len(b) > max {
		return nil, contracts.NewError(contracts.PayloadTooLarge)
	}
	out := make([]byte, 4+len(b))
	binary.BigEndian.PutUint32(out, uint32(len(b)))
	copy(out[4:], b)
	return out, nil
}
func writeFull(w io.Writer, b []byte) error {
	for len(b) > 0 {
		n, e := w.Write(b)
		if e != nil {
			return e
		}
		if n <= 0 {
			return io.ErrShortWrite
		}
		b = b[n:]
	}
	return nil
}

func safeMarshal(v any) (b []byte, err error) {
	defer func() {
		if recover() != nil {
			b = nil
			err = contracts.NewError(contracts.InternalError)
		}
	}()
	return json.Marshal(v)
}
