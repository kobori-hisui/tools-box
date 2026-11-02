// Package sidecar runs the Rust sidecar as a child process and forwards
// protocol requests to it over stdin/stdout using the same framing as the
// browser connection. Requests are multiplexed by id.
package sidecar

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"example.com/toolbox/native/internal/frame"
	"example.com/toolbox/native/pkg/protocol"
)

// ErrNotRunning is returned when the sidecar has exited or never started.
var ErrNotRunning = errors.New("sidecar is not running")

// Client is a started sidecar process.
type Client struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser

	wmu sync.Mutex // serialises frames written to stdin

	mu      sync.Mutex
	nextID  uint64
	pending map[string]chan protocol.Response
	dead    bool
	done    chan struct{}
}

// Start launches the binary at path.
func Start(path string) (*Client, error) {
	cmd := exec.Command(path)
	cmd.Stderr = os.Stderr // never stdout: that is the data channel
	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start sidecar: %w", err)
	}
	c := &Client{cmd: cmd, stdin: in, pending: map[string]chan protocol.Response{}, done: make(chan struct{})}
	go c.readLoop(out)
	return c, nil
}

func (c *Client) readLoop(r io.Reader) {
	defer close(c.done)
	for {
		buf, err := frame.ReadMax(r, frame.MaxInbound)
		if err != nil {
			break
		}
		var resp protocol.Response
		if json.Unmarshal(buf, &resp) != nil || resp.ID == "" {
			continue
		}
		c.mu.Lock()
		ch := c.pending[resp.ID]
		delete(c.pending, resp.ID)
		c.mu.Unlock()
		if ch != nil {
			ch <- resp
		}
	}
	c.mu.Lock()
	c.dead = true
	for id, ch := range c.pending {
		ch <- protocol.Response{ID: id, Error: &protocol.ErrorBody{Code: protocol.CodeUnavailable, Message: "sidecar exited"}}
		delete(c.pending, id)
	}
	c.mu.Unlock()
}

// Call forwards one request and waits for its response or ctx cancellation.
func (c *Client) Call(ctx context.Context, method string, params json.RawMessage) (json.RawMessage, error) {
	c.mu.Lock()
	if c.dead {
		c.mu.Unlock()
		return nil, ErrNotRunning
	}
	c.nextID++
	id := fmt.Sprintf("s%d", c.nextID)
	ch := make(chan protocol.Response, 1)
	c.pending[id] = ch
	c.mu.Unlock()

	payload, err := json.Marshal(protocol.Request{ID: id, Method: method, Params: params})
	if err != nil {
		return nil, err
	}
	c.wmu.Lock()
	err = frame.WriteMax(c.stdin, payload, frame.MaxInbound)
	c.wmu.Unlock()
	if err != nil {
		c.forget(id)
		return nil, err
	}

	select {
	case resp := <-ch:
		if !resp.OK {
			if resp.Error != nil {
				return nil, &RemoteError{Code: resp.Error.Code, Message: resp.Error.Message}
			}
			return nil, errors.New("sidecar returned an error")
		}
		return resp.Result, nil
	case <-ctx.Done():
		c.forget(id)
		return nil, ctx.Err()
	}
}

func (c *Client) forget(id string) {
	c.mu.Lock()
	delete(c.pending, id)
	c.mu.Unlock()
}

// Close stops the sidecar.
func (c *Client) Close() error {
	_ = c.stdin.Close()
	select {
	case <-c.done:
	default:
		_ = c.cmd.Process.Kill()
		<-c.done
	}
	return c.cmd.Wait()
}

// RemoteError is an error reported by the sidecar itself.
type RemoteError struct{ Code, Message string }

func (e *RemoteError) Error() string { return e.Code + ": " + e.Message }
