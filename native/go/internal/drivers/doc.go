// Package drivers groups low-level device access exposed to the extension.
//
// Add one sub-package per device family (serial, hid, usb), implement it behind
// a build tag when it needs cgo or OS-specific code, and register its methods
// in internal/app as "drivers.<family>.<action>". Anything that needs raw
// speed or unsafe memory access belongs in native/rust/crates/drivers and is
// reached through the "sidecar.*" prefix instead.
package drivers
