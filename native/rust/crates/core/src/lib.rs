//! Pure, dependency-light logic shared by the WASM module and the sidecar.
//! Keep it free of I/O so it compiles for `wasm32-unknown-unknown`.

pub mod search;
