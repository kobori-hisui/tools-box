package installer

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestChromeManifest(t *testing.T) {
	m := Chrome("/opt/toolbox/toolbox-host", []string{"abcdefghijklmnopabcdefghijklmnop"})
	b, err := Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got["type"] != "stdio" || got["name"] != HostName {
		t.Fatalf("bad manifest: %s", b)
	}
	if _, ok := got["allowed_extensions"]; ok {
		t.Fatal("chrome manifest must not contain allowed_extensions")
	}
	origins := got["allowed_origins"].([]any)
	if origins[0] != "chrome-extension://abcdefghijklmnopabcdefghijklmnop/" {
		t.Fatalf("origin must end with a slash: %v", origins)
	}
}

func TestFirefoxManifest(t *testing.T) {
	b, _ := Marshal(Firefox("/x/host", []string{"toolbox@example.com"}))
	s := string(b)
	if !strings.Contains(s, `"allowed_extensions"`) || strings.Contains(s, "allowed_origins") {
		t.Fatalf("bad firefox manifest: %s", s)
	}
}

func TestTargets(t *testing.T) {
	for _, goos := range []string{"linux", "darwin", "windows"} {
		ts := Targets(goos, "/home/u")
		if len(ts) != 4 {
			t.Fatalf("%s: want 4 targets, got %d", goos, len(ts))
		}
		for _, tg := range ts {
			if goos == "windows" {
				if !strings.HasSuffix(tg.RegKey, `\`+HostName) || tg.Dir != "" {
					t.Errorf("windows target %+v", tg)
				}
			} else if filepath.Base(tg.File()) != HostName+".json" || tg.RegKey != "" {
				t.Errorf("%s target %+v", goos, tg)
			}
		}
	}
	if Targets("plan9", "/h") != nil {
		t.Fatal("unknown OS must have no targets")
	}
}

func TestValidate(t *testing.T) {
	if Validate(Chrome("relative/host", []string{"id"})) == nil {
		t.Error("relative path must fail")
	}
	if Validate(Chrome("/abs/host", nil)) == nil {
		t.Error("missing ids must fail")
	}
	if err := Validate(Chrome("/abs/host", []string{"id"})); err != nil {
		t.Error(err)
	}
}
