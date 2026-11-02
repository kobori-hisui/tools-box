package protocol

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func load(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "protocol", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// strict decodes with DisallowUnknownFields, re-encodes and compares as JSON
// values, so any drift between the fixtures and the Go types fails the test.
func strict[T any](t *testing.T, name string) {
	t.Helper()
	raw := load(t, name)
	var v T
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	re, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var a, b any
	_ = json.Unmarshal(raw, &a)
	_ = json.Unmarshal(re, &b)
	ra, _ := json.Marshal(a)
	rb, _ := json.Marshal(b)
	if !bytes.Equal(ra, rb) {
		t.Fatalf("%s changed after round trip:\n%s\n%s", name, ra, rb)
	}
}

func TestFixtures(t *testing.T) {
	strict[Request](t, "request.ping.json")
	strict[Request](t, "request.fs-watch.json")
	strict[Response](t, "response.ping.ok.json")
	strict[Response](t, "response.error.json")
	strict[Event](t, "event.fs-changed.json")
}
