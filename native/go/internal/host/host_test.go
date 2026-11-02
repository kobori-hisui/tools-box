package host

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"example.com/toolbox/native/internal/app"
	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/internal/frame"
	"example.com/toolbox/native/pkg/protocol"
)

// fixture reads a golden file from the repo-level /protocol/fixtures.
func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "protocol", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return bytes.TrimSpace(b)
}

func run(t *testing.T, cfg config.Config, requests ...[]byte) []protocol.Response {
	t.Helper()
	var in, out bytes.Buffer
	for _, r := range requests {
		if err := frame.Write(&in, r); err != nil {
			t.Fatal(err)
		}
	}
	srv := New(&out)
	a, err := app.New(cfg, srv.Emit)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	srv.SetRouter(a.Router)
	if err := srv.Serve(context.Background(), &in); err != nil {
		t.Fatal(err)
	}
	var resps []protocol.Response
	for {
		b, err := frame.Read(&out)
		if err != nil {
			break
		}
		var r protocol.Response
		if err := json.Unmarshal(b, &r); err != nil {
			t.Fatal(err)
		}
		resps = append(resps, r)
	}
	return resps
}

func TestPingFixture(t *testing.T) {
	resps := run(t, config.Config{}, fixture(t, "request.ping.json"))
	if len(resps) != 1 || !resps[0].OK || resps[0].ID != "1" {
		t.Fatalf("unexpected: %+v", resps)
	}
	var got, want struct {
		Pong    bool   `json:"pong"`
		Version string `json:"version"`
	}
	_ = json.Unmarshal(resps[0].Result, &got)
	var golden protocol.Response
	_ = json.Unmarshal(fixture(t, "response.ping.ok.json"), &golden)
	_ = json.Unmarshal(golden.Result, &want)
	if !got.Pong || got.Version != want.Version {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestFsWatchForbiddenWithoutRoots(t *testing.T) {
	resps := run(t, config.Config{}, fixture(t, "request.fs-watch.json"))
	if len(resps) != 1 || resps[0].OK || resps[0].Error.Code != protocol.CodeForbidden {
		t.Fatalf("expected forbidden, got %+v", resps)
	}
}

func TestBadJSONAndUnknownMethod(t *testing.T) {
	resps := run(t, config.Config{}, []byte(`{not json`), []byte(`{"id":"9","method":"does.not.exist"}`))
	codes := map[string]string{}
	for _, r := range resps {
		codes[r.ID] = r.Error.Code
	}
	if codes[""] != protocol.CodeBadRequest || codes["9"] != protocol.CodeMethodNotFound {
		t.Fatalf("got %v", codes)
	}
}

func TestSidecarUnavailable(t *testing.T) {
	resps := run(t, config.Config{}, []byte(`{"id":"5","method":"sidecar.ping"}`))
	if len(resps) != 1 || resps[0].Error == nil || resps[0].Error.Code != protocol.CodeUnavailable {
		t.Fatalf("got %+v", resps)
	}
}

func TestOversizeResponseBecomesTooLarge(t *testing.T) {
	var out bytes.Buffer
	srv := New(&out)
	srv.respond(protocol.Response{ID: "7", OK: true, Result: json.RawMessage(`"` + strings.Repeat("x", frame.MaxOutbound) + `"`)})
	b, err := frame.Read(&out)
	if err != nil {
		t.Fatal(err)
	}
	var r protocol.Response
	_ = json.Unmarshal(b, &r)
	if r.ID != "7" || r.Error == nil || r.Error.Code != protocol.CodeTooLarge {
		t.Fatalf("got %+v", r)
	}
}
