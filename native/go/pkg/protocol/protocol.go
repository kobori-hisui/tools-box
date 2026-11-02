// Package protocol defines the wire format shared by the extension (TypeScript),
// the Go host/daemon and the Rust sidecar. The source of truth is the set of
// golden files in /protocol/fixtures, which every language tests against.
package protocol

import "encoding/json"

// Version of the protocol / binaries. Bump together with the extension.
const Version = "0.1.0"

// Request is sent by the extension. Params is optional.
type Request struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// ErrorBody describes a failed request.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response answers exactly one Request (same ID).
type Response struct {
	ID     string          `json:"id"`
	OK     bool            `json:"ok"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ErrorBody      `json:"error,omitempty"`
}

// Event is pushed by the host without a request (no ID).
type Event struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data,omitempty"`
}

// Error codes. Keep in sync with src/core/native/protocol.ts.
const (
	CodeBadRequest     = "bad_request"
	CodeMethodNotFound = "method_not_found"
	CodeForbidden      = "forbidden"
	CodeTooLarge       = "too_large"
	CodeTimeout        = "timeout"
	CodeInternal       = "internal"
	CodeUnavailable    = "unavailable"
)
