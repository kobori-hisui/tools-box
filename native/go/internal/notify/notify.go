// Package notify shows a desktop notification using the platform's own tool.
package notify

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ErrUnsupported is returned when no notification mechanism is available.
var ErrUnsupported = errors.New("notifications are not supported on this platform")

// Show displays a notification. Arguments are passed as argv (never through a
// shell) so titles and bodies cannot inject commands.
func Show(ctx context.Context, title, body string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd":
		if _, err := exec.LookPath("notify-send"); err != nil {
			return fmt.Errorf("%w: notify-send not found", ErrUnsupported)
		}
		cmd = exec.CommandContext(ctx, "notify-send", "--", title, body)
	case "darwin":
		script := fmt.Sprintf("display notification %s with title %s", appleString(body), appleString(title))
		cmd = exec.CommandContext(ctx, "osascript", "-e", script)
	default:
		// Windows toast notifications need a registered AUMID (COM/WinRT);
		// plug a dedicated implementation in here (see docs/NATIVE.md).
		return fmt.Errorf("%w: %s", ErrUnsupported, runtime.GOOS)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("notify: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// appleString quotes s as an AppleScript string literal.
func appleString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}
