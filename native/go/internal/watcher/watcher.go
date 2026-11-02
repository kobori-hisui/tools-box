// Package watcher reports filesystem changes.
//
// This implementation polls with the standard library only. It is correct
// everywhere but not efficient for huge trees. To use kernel notifications
// (inotify / FSEvents / ReadDirectoryChangesW) implement Backend with
// github.com/fsnotify/fsnotify and pass it to NewManager - nothing else changes.
package watcher

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
	"time"
)

// Event is one observed change. Op is "create", "write" or "remove".
type Event struct {
	WatchID string `json:"watchId"`
	Path    string `json:"path"`
	Op      string `json:"op"`
}

// Backend watches one path and calls emit for each change until stop is closed.
type Backend interface {
	Watch(path string, recursive bool, emit func(op, path string), stop <-chan struct{}) error
}

// Manager owns the active watches.
type Manager struct {
	backend Backend
	emit    func(Event)

	mu      sync.Mutex
	next    int
	watches map[string]chan struct{}
}

// NewManager builds a manager using the given backend.
func NewManager(b Backend, emit func(Event)) *Manager {
	return &Manager{backend: b, emit: emit, watches: map[string]chan struct{}{}}
}

// NewPolling builds a manager that polls every interval.
func NewPolling(interval time.Duration, emit func(Event)) *Manager {
	return NewManager(&Poller{Interval: interval}, emit)
}

// Add starts watching path and returns the watch id.
func (m *Manager) Add(path string, recursive bool) (string, error) {
	m.mu.Lock()
	m.next++
	id := fmt.Sprintf("w%d", m.next)
	stop := make(chan struct{})
	m.watches[id] = stop
	m.mu.Unlock()

	ready := make(chan error, 1)
	go func() {
		err := m.backend.Watch(path, recursive, func(op, p string) {
			m.emit(Event{WatchID: id, Path: p, Op: op})
		}, stop)
		select {
		case ready <- err:
		default:
		}
	}()
	// Surface immediate failures (e.g. path does not exist).
	select {
	case err := <-ready:
		if err != nil {
			m.Remove(id)
			return "", err
		}
	case <-time.After(50 * time.Millisecond):
	}
	return id, nil
}

// Remove stops a watch. It reports whether the id existed.
func (m *Manager) Remove(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	stop, ok := m.watches[id]
	if ok {
		close(stop)
		delete(m.watches, id)
	}
	return ok
}

// Close stops every watch.
func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, stop := range m.watches {
		close(stop)
		delete(m.watches, id)
	}
}

// Poller is the stdlib polling Backend.
type Poller struct{ Interval time.Duration }

type state struct {
	mod  time.Time
	size int64
}

func snapshot(root string, recursive bool) (map[string]state, error) {
	out := map[string]state{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == root {
				return err
			}
			return nil // vanished while walking
		}
		if d.IsDir() && p != root && !recursive {
			return filepath.SkipDir
		}
		fi, err := d.Info()
		if err != nil {
			return nil
		}
		out[p] = state{fi.ModTime(), fi.Size()}
		return nil
	})
	return out, err
}

// Watch implements Backend.
func (p *Poller) Watch(path string, recursive bool, emit func(op, path string), stop <-chan struct{}) error {
	iv := p.Interval
	if iv <= 0 {
		iv = time.Second
	}
	prev, err := snapshot(path, recursive)
	if err != nil {
		return err
	}
	t := time.NewTicker(iv)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return nil
		case <-t.C:
		}
		cur, err := snapshot(path, recursive)
		if err != nil {
			continue
		}
		for p, s := range cur {
			old, ok := prev[p]
			switch {
			case !ok:
				emit("create", p)
			case old != s:
				emit("write", p)
			}
		}
		for p := range prev {
			if _, ok := cur[p]; !ok {
				emit("remove", p)
			}
		}
		prev = cur
	}
}
