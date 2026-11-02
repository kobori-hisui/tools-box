// Package rpc routes protocol requests to handlers and turns every outcome
// (success, typed error, panic) into a protocol.Response.
package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"example.com/toolbox/native/pkg/protocol"
)

// Handler serves one exact method.
type Handler func(ctx context.Context, params json.RawMessage) (any, error)

// MethodHandler serves every method that starts with a registered prefix and
// receives the full method name (used to forward "sidecar.*").
type MethodHandler func(ctx context.Context, method string, params json.RawMessage) (any, error)

// Error is a handler failure with a protocol error code.
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Code + ": " + e.Message }

// Errorf builds an *Error.
func Errorf(code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Fail builds a failed response.
func Fail(id, code, message string) protocol.Response {
	return protocol.Response{ID: id, Error: &protocol.ErrorBody{Code: code, Message: message}}
}

type prefixEntry struct {
	prefix string
	h      MethodHandler
}

// Router dispatches requests. It is safe for concurrent use.
type Router struct {
	mu       sync.RWMutex
	exact    map[string]Handler
	prefixes []prefixEntry
}

// NewRouter creates an empty router.
func NewRouter() *Router { return &Router{exact: map[string]Handler{}} }

// Register adds an exact-match handler.
func (r *Router) Register(method string, h Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exact[method] = h
}

// RegisterPrefix adds a prefix handler. The longest matching prefix wins.
func (r *Router) RegisterPrefix(prefix string, h MethodHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prefixes = append(r.prefixes, prefixEntry{prefix, h})
	sort.Slice(r.prefixes, func(i, j int) bool { return len(r.prefixes[i].prefix) > len(r.prefixes[j].prefix) })
}

// Methods lists exact methods and "prefix*" entries, sorted.
func (r *Router) Methods() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.exact)+len(r.prefixes))
	for m := range r.exact {
		out = append(out, m)
	}
	for _, p := range r.prefixes {
		out = append(out, p.prefix+"*")
	}
	sort.Strings(out)
	return out
}

func (r *Router) lookup(method string) func(context.Context, json.RawMessage) (any, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if h, ok := r.exact[method]; ok {
		return h
	}
	for _, p := range r.prefixes {
		if strings.HasPrefix(method, p.prefix) {
			h := p.h
			return func(ctx context.Context, params json.RawMessage) (any, error) { return h(ctx, method, params) }
		}
	}
	return nil
}

// Dispatch runs the request and always returns a response with req.ID.
func (r *Router) Dispatch(ctx context.Context, req protocol.Request) (resp protocol.Response) {
	defer func() {
		if rec := recover(); rec != nil {
			resp = Fail(req.ID, protocol.CodeInternal, fmt.Sprintf("handler panic: %v", rec))
		}
	}()
	if req.Method == "" {
		return Fail(req.ID, protocol.CodeBadRequest, "missing method")
	}
	call := r.lookup(req.Method)
	if call == nil {
		return Fail(req.ID, protocol.CodeMethodNotFound, fmt.Sprintf("unknown method %q", req.Method))
	}
	result, err := call(ctx, req.Params)
	if err != nil {
		var re *Error
		if errors.As(err, &re) {
			return Fail(req.ID, re.Code, re.Message)
		}
		return Fail(req.ID, protocol.CodeInternal, err.Error())
	}
	resp = protocol.Response{ID: req.ID, OK: true}
	if rm, ok := result.(json.RawMessage); ok {
		if len(rm) > 0 {
			resp.Result = rm
		}
		return resp
	}
	if result != nil {
		raw, err := json.Marshal(result)
		if err != nil {
			return Fail(req.ID, protocol.CodeInternal, "encode result: "+err.Error())
		}
		resp.Result = raw
	}
	return resp
}
