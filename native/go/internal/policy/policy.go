// Package policy decides which filesystem paths the extension may touch.
// Everything is denied unless it lies inside an explicitly allowed root.
package policy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrDenied is returned for paths outside the allowed roots.
var ErrDenied = errors.New("path is outside the allowed roots")

// Policy is an immutable allow-list of absolute directory roots.
type Policy struct{ roots []string }

// New normalises roots (absolute, symlinks resolved when they exist).
func New(roots []string) (*Policy, error) {
	p := &Policy{}
	for _, r := range roots {
		if r == "" {
			continue
		}
		abs, err := filepath.Abs(r)
		if err != nil {
			return nil, err
		}
		if real, err := filepath.EvalSymlinks(abs); err == nil {
			abs = real
		}
		p.roots = append(p.roots, filepath.Clean(abs))
	}
	return p, nil
}

// Roots returns a copy of the allowed roots.
func (p *Policy) Roots() []string { return append([]string(nil), p.roots...) }

// Check returns the cleaned absolute path if it is inside an allowed root.
// Symlinks in the existing part of the path are resolved first so a link
// inside a root cannot point outside it.
func (p *Policy) Check(path string) (string, error) {
	if len(p.roots) == 0 {
		return "", fmt.Errorf("%w: no roots configured", ErrDenied)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs = resolveExisting(filepath.Clean(abs))
	for _, root := range p.roots {
		if abs == root || strings.HasPrefix(abs, root+string(filepath.Separator)) {
			return abs, nil
		}
	}
	return "", ErrDenied
}

// resolveExisting resolves symlinks for the longest existing prefix.
func resolveExisting(path string) string {
	rest := ""
	cur := path
	for {
		if real, err := filepath.EvalSymlinks(cur); err == nil {
			return filepath.Join(real, rest)
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return path
		}
		rest = filepath.Join(filepath.Base(cur), rest)
		cur = parent
		if _, err := os.Lstat(cur); err == nil {
			if real, err := filepath.EvalSymlinks(cur); err == nil {
				return filepath.Join(real, rest)
			}
		}
	}
}
