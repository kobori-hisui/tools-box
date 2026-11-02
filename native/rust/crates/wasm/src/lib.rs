//! WASM module with a plain C ABI (no wasm-bindgen), loadable with
//! `WebAssembly.instantiate` from any extension context.
//!
//! Calling convention (all pointers are offsets in the module's linear memory):
//!   1. `tb_alloc(len)` -> ptr, copy the UTF-8 JSON input there
//!   2. `tb_search(ptr, len)` -> packed u64: (out_ptr << 32) | out_len
//!      (0 means error); input: {"tools":[...],"query":"..","limit":N}
//!   3. read out_len bytes at out_ptr (UTF-8 JSON array of ids)
//!   4. `tb_free(ptr, len)` for the input and the output buffers
//!
//! Build: cargo build --profile wasm -p toolbox-wasm --target wasm32-unknown-unknown

use serde::Deserialize;
use toolbox_core::search::{search, Tool};

#[derive(Deserialize)]
struct SearchInput {
    tools: Vec<Tool>,
    #[serde(default)]
    query: String,
    #[serde(default)]
    limit: usize,
}

/// Safe core used by the exported function and by the tests.
pub fn run_search(input: &[u8]) -> Option<Vec<u8>> {
    let inp: SearchInput = serde_json::from_slice(input).ok()?;
    serde_json::to_vec(&search(&inp.tools, &inp.query, inp.limit)).ok()
}

/// Every buffer handed to JS is a `Box<[u8]>` of exactly `len.max(1)` bytes,
/// so `tb_free` can always rebuild the same box.
fn leak(buf: Box<[u8]>) -> *mut u8 {
    Box::into_raw(buf) as *mut u8
}

#[no_mangle]
pub extern "C" fn tb_alloc(len: usize) -> *mut u8 {
    leak(vec![0u8; len.max(1)].into_boxed_slice())
}

/// # Safety
/// `ptr`/`len` must come from `tb_alloc(len)` or be the (ptr, len) returned by
/// `tb_search`, and must not be freed twice.
#[no_mangle]
pub unsafe extern "C" fn tb_free(ptr: *mut u8, len: usize) {
    if !ptr.is_null() {
        let slice = std::ptr::slice_from_raw_parts_mut(ptr, len.max(1));
        drop(Box::from_raw(slice));
    }
}

/// # Safety
/// `ptr` must point to `len` readable bytes.
#[no_mangle]
pub unsafe extern "C" fn tb_search(ptr: *const u8, len: usize) -> u64 {
    let input = std::slice::from_raw_parts(ptr, len);
    match run_search(input) {
        // Output is a JSON array, so it is never empty ("[]" at minimum).
        Some(out) => {
            let out_len = out.len();
            let p = leak(out.into_boxed_slice());
            ((p as usize as u64) << 32) | out_len as u64
        }
        None => 0,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn run_search_json_in_out() {
        let input = br#"{"tools":[{"id":"a","name":"Alpha Tool","keywords":["x"]},{"id":"b","name":"Beta","keywords":[]}],"query":"alpha"}"#;
        let out = run_search(input).unwrap();
        assert_eq!(String::from_utf8(out).unwrap(), r#"["a"]"#);
    }

    #[test]
    fn invalid_input_is_none() {
        assert!(run_search(b"not json").is_none());
    }
}
