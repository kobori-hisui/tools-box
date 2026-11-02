// Package host serves the native-messaging stdio protocol: framed JSON
// requests in on stdin, framed responses and events out on stdout.
// stdout carries ONLY frames; all logging must go to stderr.
package host

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"sync"
	"time"

	"example.com/toolbox/native/internal/frame"
	"example.com/toolbox/native/internal/rpc"
	"example.com/toolbox/native/pkg/protocol"
)

// RequestTimeout bounds every handler.
const RequestTimeout = 30 * time.Second

// Server writes frames (thread-safe) and dispatches requests.
type Server struct {
	out    io.Writer
	wmu    sync.Mutex
	router *rpc.Router
}

// New creates a server writing to out. Call SetRouter before Serve.
func New(out io.Writer) *Server { return &Server{out: out} }

// SetRouter installs the dispatcher.
func (s *Server) SetRouter(r *rpc.Router) { s.router = r }

func (s *Server) write(v any) error {
	payload, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.wmu.Lock()
	defer s.wmu.Unlock()
	return frame.Write(s.out, payload)
}

// Emit pushes an event. Events that exceed the 1 MB limit are dropped.
func (s *Server) Emit(ev protocol.Event) {
	if err := s.write(ev); err != nil {
		log.Printf("emit %s: %v", ev.Event, err)
	}
}

func (s *Server) respond(resp protocol.Response) {
	err := s.write(resp)
	if errors.Is(err, frame.ErrTooLarge) {
		// The browser would kill the connection on an oversize frame.
		err = s.write(rpc.Fail(resp.ID, protocol.CodeTooLarge, "response exceeds the 1 MB native messaging limit"))
	}
	if err != nil {
		log.Printf("respond %s: %v", resp.ID, err)
	}
}

// Serve reads requests until in is closed, then waits for in-flight handlers.
func (s *Server) Serve(ctx context.Context, in io.Reader) error {
	var wg sync.WaitGroup
	defer wg.Wait()
	for {
		buf, err := frame.Read(in)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		var req protocol.Request
		if err := json.Unmarshal(buf, &req); err != nil {
			s.respond(rpc.Fail("", protocol.CodeBadRequest, "invalid JSON: "+err.Error()))
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			cctx, cancel := context.WithTimeout(ctx, RequestTimeout)
			defer cancel()
			s.respond(s.router.Dispatch(cctx, req))
		}()
	}
}
