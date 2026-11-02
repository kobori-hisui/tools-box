//! Low-level device access (serial, HID, USB, raw sockets, memory-mapped I/O).
//!
//! Kept in its own crate so the heavy / platform-specific dependencies never
//! reach the WASM build or the sidecar's default features. Expose each device
//! family behind a Cargo feature and a `sidecar.<family>.*` method.

/// Device families this crate is meant to host.
pub const FAMILIES: &[&str] = &["serial", "hid", "usb"];

#[cfg(test)]
mod tests {
    #[test]
    fn families_are_declared() {
        assert!(super::FAMILIES.contains(&"usb"));
    }
}
