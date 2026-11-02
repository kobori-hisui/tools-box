package app

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/pkg/protocol"
)

// Talks to the real Rust sidecar. Build it first:
//
//	cargo build --release -p toolbox-sidecar   (in native/rust)
//
// The test is skipped when the binary is missing.
func TestRustSidecarEndToEnd(t *testing.T) {
	bin := filepath.Join("..", "..", "..", "rust", "target", "release", "toolbox-sidecar")
	if _, err := os.Stat(bin); err != nil {
		t.Skip("sidecar not built: " + bin)
	}
	abs, _ := filepath.Abs(bin)
	a, err := New(config.Config{Sidecar: abs}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp := a.Router.Dispatch(ctx, protocol.Request{ID: "1", Method: "sidecar.ping"})
	if !resp.OK {
		t.Fatalf("ping: %+v", resp.Error)
	}

	params := json.RawMessage(`{"tools":[
		{"id":"b64","name":"Base64 Encoder","keywords":["encode"]},
		{"id":"json","name":"JSON Formatter","keywords":["format"]}],"query":"encode"}`)
	resp = a.Router.Dispatch(ctx, protocol.Request{ID: "2", Method: "sidecar.search", Params: params})
	if !resp.OK || string(resp.Result) != `["b64"]` {
		t.Fatalf("search: ok=%v result=%s err=%+v", resp.OK, resp.Result, resp.Error)
	}

	resp = a.Router.Dispatch(ctx, protocol.Request{ID: "3", Method: "sidecar.nope"})
	if resp.OK || resp.Error.Code != protocol.CodeMethodNotFound {
		t.Fatalf("unknown method must pass the sidecar's error through: %+v", resp)
	}
}
