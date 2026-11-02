// Package installer builds and registers the native-messaging host manifests.
//
// Chrome-family browsers authorise callers with "allowed_origins"
// (chrome-extension://<id>/); Firefox uses "allowed_extensions" (add-on IDs).
// Firefox's Windows registry location is documented on MDN; the Chromium
// directories/keys below follow Chrome's documented layout - verify them with
// `toolbox-installer print` against the browsers you ship to.
package installer

import (
	"encoding/json"
	"fmt"
	"path/filepath"
)

// HostName is the name passed to connectNative()/sendNativeMessage().
// Native messaging names allow only lowercase letters, digits, '_' and '.'.
const HostName = "com.toolbox.host"

// Manifest is the JSON file the browser reads.
type Manifest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Path              string   `json:"path"`
	Type              string   `json:"type"`
	AllowedOrigins    []string `json:"allowed_origins,omitempty"`
	AllowedExtensions []string `json:"allowed_extensions,omitempty"`
}

// Chrome builds the manifest for Chromium-based browsers.
func Chrome(hostPath string, extensionIDs []string) Manifest {
	origins := make([]string, 0, len(extensionIDs))
	for _, id := range extensionIDs {
		origins = append(origins, "chrome-extension://"+id+"/")
	}
	return Manifest{Name: HostName, Description: "Toolbox native host", Path: hostPath, Type: "stdio", AllowedOrigins: origins}
}

// Firefox builds the manifest for Firefox.
func Firefox(hostPath string, addonIDs []string) Manifest {
	return Manifest{Name: HostName, Description: "Toolbox native host", Path: hostPath, Type: "stdio", AllowedExtensions: addonIDs}
}

// Marshal renders a manifest as indented JSON.
func Marshal(m Manifest) ([]byte, error) { return json.MarshalIndent(m, "", "  ") }

// Target says where one browser looks for the manifest.
type Target struct {
	Browser string // chrome, chromium, edge, firefox
	Family  string // "chromium" or "firefox" (selects the manifest flavour)
	Dir     string // file-based platforms: directory holding <HostName>.json
	RegKey  string // Windows: registry key whose default value is the manifest path
}

// File returns the manifest path for file-based platforms.
func (t Target) File() string { return filepath.Join(t.Dir, HostName+".json") }

// Targets lists the per-user locations for goos. home is the user's home dir;
// roaming is %APPDATA% on Windows (where manifests are stored).
func Targets(goos, home string) []Target {
	switch goos {
	case "linux":
		return []Target{
			{Browser: "chrome", Family: "chromium", Dir: filepath.Join(home, ".config/google-chrome/NativeMessagingHosts")},
			{Browser: "chromium", Family: "chromium", Dir: filepath.Join(home, ".config/chromium/NativeMessagingHosts")},
			{Browser: "edge", Family: "chromium", Dir: filepath.Join(home, ".config/microsoft-edge/NativeMessagingHosts")},
			{Browser: "firefox", Family: "firefox", Dir: filepath.Join(home, ".mozilla/native-messaging-hosts")},
		}
	case "darwin":
		app := filepath.Join(home, "Library/Application Support")
		return []Target{
			{Browser: "chrome", Family: "chromium", Dir: filepath.Join(app, "Google/Chrome/NativeMessagingHosts")},
			{Browser: "chromium", Family: "chromium", Dir: filepath.Join(app, "Chromium/NativeMessagingHosts")},
			{Browser: "edge", Family: "chromium", Dir: filepath.Join(app, "Microsoft Edge/NativeMessagingHosts")},
			{Browser: "firefox", Family: "firefox", Dir: filepath.Join(app, "Mozilla/NativeMessagingHosts")},
		}
	case "windows":
		return []Target{
			{Browser: "chrome", Family: "chromium", RegKey: `HKCU\Software\Google\Chrome\NativeMessagingHosts\` + HostName},
			{Browser: "chromium", Family: "chromium", RegKey: `HKCU\Software\Chromium\NativeMessagingHosts\` + HostName},
			{Browser: "edge", Family: "chromium", RegKey: `HKCU\Software\Microsoft\Edge\NativeMessagingHosts\` + HostName},
			{Browser: "firefox", Family: "firefox", RegKey: `HKCU\Software\Mozilla\NativeMessagingHosts\` + HostName},
		}
	}
	return nil
}

// ManifestFor picks the right flavour for a target.
func ManifestFor(t Target, hostPath string, chromeIDs, firefoxIDs []string) Manifest {
	if t.Family == "firefox" {
		return Firefox(hostPath, firefoxIDs)
	}
	return Chrome(hostPath, chromeIDs)
}

// Validate rejects manifests the browser would refuse.
func Validate(m Manifest) error {
	if !filepath.IsAbs(m.Path) {
		return fmt.Errorf("host path must be absolute: %q", m.Path)
	}
	if len(m.AllowedOrigins) == 0 && len(m.AllowedExtensions) == 0 {
		return fmt.Errorf("no allowed extension ids for %s", m.Name)
	}
	return nil
}
