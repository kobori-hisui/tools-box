package policy

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCheck(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := New([]string{root})
	if err != nil {
		t.Fatal(err)
	}

	ok := []string{root, filepath.Join(root, "sub"), filepath.Join(root, "sub", "new.txt"), filepath.Join(root, "sub", "..", "sub")}
	for _, path := range ok {
		if _, err := p.Check(path); err != nil {
			t.Errorf("Check(%q) denied: %v", path, err)
		}
	}
	bad := []string{outside, filepath.Join(root, ".."), filepath.Join(root+"-evil", "x"), "/etc/passwd"}
	for _, path := range bad {
		if _, err := p.Check(path); !errors.Is(err, ErrDenied) {
			t.Errorf("Check(%q) = %v, want ErrDenied", path, err)
		}
	}
}

func TestSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks need privileges on Windows")
	}
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "link")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	p, _ := New([]string{root})
	if _, err := p.Check(filepath.Join(link, "secret.txt")); !errors.Is(err, ErrDenied) {
		t.Fatalf("symlink escape allowed: %v", err)
	}
}

func TestNoRootsDeniesAll(t *testing.T) {
	p, _ := New(nil)
	if _, err := p.Check(t.TempDir()); !errors.Is(err, ErrDenied) {
		t.Fatalf("want ErrDenied, got %v", err)
	}
}
