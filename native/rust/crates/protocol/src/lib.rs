//! Wire format shared with the extension and the Go host.
//! The golden files in `/protocol/fixtures` are the source of truth.

use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::io::{self, Read, Write};

pub const VERSION: &str = "0.1.0";
/// Largest frame the sidecar accepts from the host.
pub const MAX_FRAME: usize = 64 << 20;

pub mod code {
    pub const BAD_REQUEST: &str = "bad_request";
    pub const METHOD_NOT_FOUND: &str = "method_not_found";
    pub const FORBIDDEN: &str = "forbidden";
    pub const TOO_LARGE: &str = "too_large";
    pub const TIMEOUT: &str = "timeout";
    pub const INTERNAL: &str = "internal";
    pub const UNAVAILABLE: &str = "unavailable";
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Request {
    pub id: String,
    pub method: String,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub params: Option<Value>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct ErrorBody {
    pub code: String,
    pub message: String,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Response {
    pub id: String,
    pub ok: bool,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub result: Option<Value>,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub error: Option<ErrorBody>,
}

#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(deny_unknown_fields)]
pub struct Event {
    pub event: String,
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub data: Option<Value>,
}

impl Response {
    pub fn ok(id: impl Into<String>, result: Value) -> Self {
        Self {
            id: id.into(),
            ok: true,
            result: Some(result),
            error: None,
        }
    }
    pub fn fail(id: impl Into<String>, code: &str, message: impl Into<String>) -> Self {
        Self {
            id: id.into(),
            ok: false,
            result: None,
            error: Some(ErrorBody {
                code: code.into(),
                message: message.into(),
            }),
        }
    }
}

/// Read one frame. `Ok(None)` means clean EOF between frames.
pub fn read_frame<R: Read>(r: &mut R) -> io::Result<Option<Vec<u8>>> {
    let mut hdr = [0u8; 4];
    match r.read_exact(&mut hdr) {
        Ok(()) => {}
        Err(e) if e.kind() == io::ErrorKind::UnexpectedEof => {
            // EOF before any header byte is clean; we cannot tell partial
            // headers apart with read_exact, which is acceptable for a pipe.
            return Ok(None);
        }
        Err(e) => return Err(e),
    }
    let n = u32::from_ne_bytes(hdr) as usize;
    if n > MAX_FRAME {
        return Err(io::Error::new(
            io::ErrorKind::InvalidData,
            format!("frame of {n} bytes exceeds {MAX_FRAME}"),
        ));
    }
    let mut buf = vec![0u8; n];
    r.read_exact(&mut buf)?;
    Ok(Some(buf))
}

/// Write one frame (native byte order length prefix) in a single write call.
pub fn write_frame<W: Write>(w: &mut W, payload: &[u8]) -> io::Result<()> {
    let len = u32::try_from(payload.len())
        .map_err(|_| io::Error::new(io::ErrorKind::InvalidInput, "payload larger than u32::MAX"))?;
    let mut buf = Vec::with_capacity(4 + payload.len());
    buf.extend_from_slice(&len.to_ne_bytes());
    buf.extend_from_slice(payload);
    w.write_all(&buf)?;
    w.flush()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn frame_round_trip() {
        let mut buf = Vec::new();
        write_frame(&mut buf, b"{\"a\":1}").unwrap();
        write_frame(&mut buf, b"").unwrap();
        let mut cur = io::Cursor::new(buf);
        assert_eq!(read_frame(&mut cur).unwrap().unwrap(), b"{\"a\":1}");
        assert_eq!(read_frame(&mut cur).unwrap().unwrap(), b"");
        assert!(read_frame(&mut cur).unwrap().is_none());
    }

    #[test]
    fn oversize_header_is_rejected() {
        let hdr = ((MAX_FRAME + 1) as u32).to_ne_bytes();
        assert!(read_frame(&mut io::Cursor::new(hdr.to_vec())).is_err());
    }

    #[test]
    fn truncated_payload_is_error() {
        let mut buf = Vec::new();
        write_frame(&mut buf, b"hello").unwrap();
        buf.truncate(6);
        assert!(read_frame(&mut io::Cursor::new(buf)).is_err());
    }
}
