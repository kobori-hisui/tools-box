package daemon

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/toolbox/native/internal/app"
	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/pkg/protocol"
)

const token = "0123456789abcdef-token"

func newServer(t *testing.T) (*httptest.Server, *Hub) {
	t.Helper()
	hub := NewHub()
	a, err := app.New(config.Config{}, hub.Publish)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(a.Close)
	h, err := Handler(Options{Router: a.Router, Hub: hub, Token: token, AllowedOrigins: []string{"chrome-extension://abc"}})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv, hub
}

func do(t *testing.T, method, url, body string, headers map[string]string) *http.Response {
	t.Helper()
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { res.Body.Close() })
	return res
}

func TestShortTokenRefused(t *testing.T) {
	if _, err := Handler(Options{Token: "short"}); err == nil {
		t.Fatal("expected error for short token")
	}
}

func TestAuthAndRPC(t *testing.T) {
	srv, _ := newServer(t)
	ping := `{"id":"1","method":"ping"}`

	if res := do(t, "POST", srv.URL+"/rpc", ping, nil); res.StatusCode != 401 {
		t.Fatalf("no token: %d", res.StatusCode)
	}
	if res := do(t, "POST", srv.URL+"/rpc", ping, map[string]string{"x-api-key": "wrong-wrong-wrong-wrong"}); res.StatusCode != 401 {
		t.Fatalf("bad token: %d", res.StatusCode)
	}
	res := do(t, "POST", srv.URL+"/rpc", ping, map[string]string{"x-api-key": token})
	var resp protocol.Response
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil || !resp.OK || resp.ID != "1" {
		t.Fatalf("rpc: %v %+v", err, resp)
	}
	if res := do(t, "GET", srv.URL+"/health", "", nil); res.StatusCode != 200 {
		t.Fatalf("health: %d", res.StatusCode)
	}
}

func TestOrigins(t *testing.T) {
	srv, _ := newServer(t)
	good := map[string]string{"Origin": "chrome-extension://abc", "x-api-key": token}
	res := do(t, "POST", srv.URL+"/rpc", `{"id":"1","method":"ping"}`, good)
	if res.StatusCode != 200 || res.Header.Get("Access-Control-Allow-Origin") != "chrome-extension://abc" {
		t.Fatalf("allowed origin: %d %v", res.StatusCode, res.Header)
	}
	if res := do(t, "OPTIONS", srv.URL+"/rpc", "", map[string]string{"Origin": "chrome-extension://abc"}); res.StatusCode != 204 {
		t.Fatalf("preflight: %d", res.StatusCode)
	}
	evil := map[string]string{"Origin": "https://evil.example", "x-api-key": token}
	if res := do(t, "POST", srv.URL+"/rpc", `{"id":"1","method":"ping"}`, evil); res.StatusCode != 403 {
		t.Fatalf("evil origin: %d", res.StatusCode)
	}
}

func TestEventsStream(t *testing.T) {
	srv, hub := newServer(t)
	req, _ := http.NewRequest("GET", srv.URL+"/events", nil)
	req.Header.Set("x-api-key", token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	r := bufio.NewReader(res.Body)
	if line, _ := r.ReadString('\n'); !strings.HasPrefix(line, ": connected") {
		t.Fatalf("unexpected first line %q", line)
	}
	_, _ = r.ReadString('\n')

	got := make(chan string, 1)
	go func() {
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				return
			}
			if strings.HasPrefix(line, "data: ") {
				got <- line
				return
			}
		}
	}()
	hub.Publish(protocol.Event{Event: "fs.changed", Data: json.RawMessage(`{"path":"/x"}`)})
	select {
	case line := <-got:
		if !strings.Contains(line, `"event":"fs.changed"`) {
			t.Fatalf("bad event %q", line)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no event received")
	}
}
