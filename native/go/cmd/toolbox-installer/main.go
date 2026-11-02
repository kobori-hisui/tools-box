// toolbox-installer registers (or removes) the native-messaging host manifest
// for the browsers installed on this machine.
//
//	toolbox-installer print     --host /abs/toolbox-host --chrome-id ID --firefox-id ID
//	toolbox-installer install   --host /abs/toolbox-host --chrome-id ID --firefox-id ID
//	toolbox-installer uninstall
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/internal/installer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: toolbox-installer print|install|uninstall [flags]")
		os.Exit(2)
	}
	cmd := os.Args[1]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	host := fs.String("host", "", "absolute path to the toolbox-host binary")
	chromeIDs := fs.String("chrome-id", "", "comma-separated Chrome/Edge extension ids")
	firefoxIDs := fs.String("firefox-id", "", "comma-separated Firefox add-on ids")
	_ = fs.Parse(os.Args[2:])

	home, err := os.UserHomeDir()
	check(err)
	targets := installer.Targets(runtime.GOOS, home)
	if targets == nil {
		check(fmt.Errorf("unsupported OS %q", runtime.GOOS))
	}

	switch cmd {
	case "print", "install":
		abs, err := filepath.Abs(*host)
		if *host == "" || err != nil {
			check(errors.New("--host is required"))
		}
		for _, t := range targets {
			m := installer.ManifestFor(t, abs, split(*chromeIDs), split(*firefoxIDs))
			if len(m.AllowedOrigins)+len(m.AllowedExtensions) == 0 {
				fmt.Printf("skip %-8s (no extension id given for this browser family)\n", t.Browser)
				continue
			}
			check(installer.Validate(m))
			data, err := installer.Marshal(m)
			check(err)
			if cmd == "print" {
				fmt.Printf("# %s -> %s%s\n%s\n", t.Browser, t.Dir, t.RegKey, data)
				continue
			}
			check(install(t, data))
			fmt.Printf("installed %-8s\n", t.Browser)
		}
	case "uninstall":
		for _, t := range targets {
			check(uninstall(t))
			fmt.Printf("removed %-8s\n", t.Browser)
		}
	default:
		check(fmt.Errorf("unknown command %q", cmd))
	}
}

func install(t installer.Target, data []byte) error {
	file := t.File()
	if t.RegKey != "" {
		dir, err := config.Dir()
		if err != nil {
			return err
		}
		file = filepath.Join(dir, "manifests", t.Browser+".json")
	}
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(file, data, 0o644); err != nil {
		return err
	}
	if t.RegKey != "" {
		out, err := exec.Command("reg", "add", t.RegKey, "/ve", "/t", "REG_SZ", "/d", file, "/f").CombinedOutput()
		if err != nil {
			return fmt.Errorf("reg add: %w: %s", err, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func uninstall(t installer.Target) error {
	if t.RegKey != "" {
		// Ignore "not found": uninstall must be idempotent.
		_ = exec.Command("reg", "delete", t.RegKey, "/f").Run()
		return nil
	}
	if err := os.Remove(t.File()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func split(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "toolbox-installer:", err)
		os.Exit(1)
	}
}
