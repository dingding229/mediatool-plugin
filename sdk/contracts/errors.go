// Package contracts is the public, persistence-independent plugin wire contract.
package contracts

import "fmt"

const ContractVersion = 2

type Code int

const (
	ParseError             Code = -32700
	InvalidRequest         Code = -32600
	MethodNotFound         Code = -32601
	InvalidParams          Code = -32602
	InternalError          Code = -32603
	PermissionDenied       Code = -32001
	Unavailable            Code = -32002
	DeadlineExceeded       Code = -32003
	Cancelled              Code = -32004
	Conflict               Code = -32005
	Busy                   Code = -32006
	NotFound               Code = -32007
	UnsupportedVersion     Code = -32008
	UpstreamLimited        Code = -32009
	NeedsAction            Code = -32010
	ResultUnknown          Code = -32011
	PayloadTooLarge        Code = -32012
	InvalidExtensionResult Code = -32013
	Unsupported            Code = -32014
)

type ErrorData struct {
	Reason       string `json:"reason"`
	RetryAfterMS int64  `json:"retryAfterMs,omitempty"`
}
type Error struct {
	Code    Code      `json:"code"`
	Message string    `json:"message"`
	Data    ErrorData `json:"data"`
}

func Reason(c Code) string {
	switch c {
	case Unsupported:
		return "unsupported"
	case ParseError:
		return "parse_error"
	case InvalidRequest:
		return "invalid_request"
	case MethodNotFound:
		return "method_not_found"
	case InvalidParams:
		return "invalid_params"
	case PermissionDenied:
		return "permission_denied"
	case Unavailable:
		return "unavailable"
	case DeadlineExceeded:
		return "deadline_exceeded"
	case Cancelled:
		return "cancelled"
	case Conflict:
		return "conflict"
	case Busy:
		return "busy"
	case NotFound:
		return "not_found"
	case UnsupportedVersion:
		return "unsupported_version"
	case UpstreamLimited:
		return "upstream_limited"
	case NeedsAction:
		return "needs_action"
	case ResultUnknown:
		return "result_unknown"
	case PayloadTooLarge:
		return "payload_too_large"
	case InvalidExtensionResult:
		return "invalid_extension_result"
	default:
		return "internal_error"
	}
}
func NewError(c Code) *Error {
	return &Error{Code: c, Message: Reason(c), Data: ErrorData{Reason: Reason(c)}}
}

// Error never includes remote messages or payloads, which can contain secrets.
func (e *Error) Error() string { return fmt.Sprintf("plugin rpc %d: %s", e.Code, Reason(e.Code)) }

// SafeError never relays arbitrary callback/remote error messages or data.
func SafeError(e *Error) *Error {
	if e == nil {
		return NewError(InternalError)
	}
	code := e.Code
	if Reason(code) == "internal_error" {
		code = InternalError
	}
	safe := NewError(code)
	if e.Data.RetryAfterMS > 0 && e.Data.RetryAfterMS <= 86400000 {
		safe.Data.RetryAfterMS = e.Data.RetryAfterMS
	}
	return safe
}
