// Package daemon exposes the same protocol over a loopback HTTP server for
// work that must outlive the browser (file watching, schedulers, long jobs).
//
//	GET  /health   no auth, no secrets
//	POST /rpc      body: protocol.Request      -> protocol.Response
//	GET  /events   server-sent events (protocol.Event as JSON)
//
// Every route except /health needs the x-api-key header. Browser callers must
// also come from an allowed origin (chrome-extension://<id>, moz-extension://<uuid>).
package daemon

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"example.com/toolbox/native/internal/rpc"
	"example.com/toolbox/native/pkg/protocol"
)

// MaxBody caps a request body.
const MaxBody = 8 << 20

// Hub fans events out to SSE subscribers; slow subscribers drop events.
type Hub struct {
	mu   sync.Mutex
	subs map[chan protocol.Event]struct{}
}

// NewHub creates an empty hub.
func NewHub() *Hub { return &Hub{subs: map[chan protocol.Event]struct{}{}} }

// Publish delivers ev to every subscriber without blocking.
func (h *Hub) Publish(ev protocol.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.subs {
		select {
		case ch <- ev:
		default:
		}
	}
}

func (h *Hub) subscribe() (chan protocol.Event, func()) {
	ch := make(chan protocol.Event, 64)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	return ch, func() {
		h.mu.Lock()
		delete(h.subs, ch)
		h.mu.Unlock()
	}
}

// Options configure the handler.
type Options struct {
	Router         *rpc.Router
	Hub            *Hub
	Token          string
	AllowedOrigins []string
}

// Handler builds the HTTP handler. It fails closed if no token is set.
func Handler(o Options) (http.Handler, error) {
	if len(o.Token) < 16 {
		return nil, errors.New("daemon: token must be at least 16 characters")
	}
	if o.Hub == nil {
		o.Hub = NewHub()
	}
	origins := map[string]bool{}
	for _, v := range o.AllowedOrigins {
		origins[v] = true
	}
	s := &server{o: o, origins: origins}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/rpc", s.guard(s.rpc))
	mux.HandleFunc("/events", s.guard(s.events))
	return mux, nil
}

type server struct {
	o       Options
	origins map[string]bool
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"ok":true,"version":%q}`, protocol.Version)
}

func (s *server) guard(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			if !s.origins[origin] {
				http.Error(w, "origin not allowed", http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "content-type, x-api-key")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		key := r.Header.Get("x-api-key")
		if subtle.ConstantTimeCompare([]byte(key), []byte(s.o.Token)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (s *server) rpc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req protocol.Request
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, MaxBody))
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, rpc.Fail("", protocol.CodeBadRequest, "invalid JSON: "+err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	writeJSON(w, http.StatusOK, s.o.Router.Dispatch(ctx, req))
}

func (s *server) events(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	ch, cancel := s.o.Hub.subscribe()
	defer cancel()
	fmt.Fprint(w, ": connected\n\n")
	fl.Flush()
	keepalive := time.NewTicker(20 * time.Second)
	defer keepalive.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case ev := <-ch:
			b, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", b)
			fl.Flush()
		case <-keepalive.C:
			fmt.Fprint(w, ": keepalive\n\n")
			fl.Flush()
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
