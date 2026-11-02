package watcher

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

type collector struct {
	mu  sync.Mutex
	got []Event
}

func (c *collector) add(e Event) { c.mu.Lock(); c.got = append(c.got, e); c.mu.Unlock() }
func (c *collector) has(op, path string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range c.got {
		if e.Op == op && e.Path == path {
			return true
		}
	}
	return false
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("timed out waiting for event")
}

func TestPollingCreateWriteRemove(t *testing.T) {
	dir := t.TempDir()
	c := &collector{}
	m := NewPolling(15*time.Millisecond, c.add)
	defer m.Close()
	id, err := m.Add(dir, true)
	if err != nil {
		t.Fatal(err)
	}

	f := filepath.Join(dir, "a.txt")
	if err := os.WriteFile(f, []byte("1"), 0o644); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { return c.has("create", f) })

	if err := os.WriteFile(f, []byte("changed!"), 0o644); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { return c.has("write", f) })

	if err := os.Remove(f); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { return c.has("remove", f) })

	if !m.Remove(id) || m.Remove(id) {
		t.Fatal("Remove should succeed once, then report unknown id")
	}
}

func TestAddMissingPath(t *testing.T) {
	m := NewPolling(10*time.Millisecond, func(Event) {})
	defer m.Close()
	if _, err := m.Add(filepath.Join(t.TempDir(), "nope"), false); err == nil {
		t.Fatal("expected error for missing path")
	}
}
