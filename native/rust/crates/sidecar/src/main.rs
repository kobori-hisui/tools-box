//! Heavy-work sidecar. The Go host starts it and forwards every `sidecar.*`
//! request over stdin/stdout (same framing as native messaging).
//! stdout carries ONLY frames; diagnostics go to stderr.

use serde::Deserialize;
use serde_json::json;
use std::io::{self, BufWriter};
use toolbox_core::search::{search, Tool};
use toolbox_protocol::{code, read_frame, write_frame, Request, Response, VERSION};

#[derive(Deserialize)]
struct SearchParams {
    tools: Vec<Tool>,
    #[serde(default)]
    query: String,
    #[serde(default)]
    limit: usize,
}

fn handle(req: Request) -> Response {
    let id = req.id.clone();
    match req.method.as_str() {
        "sidecar.ping" => Response::ok(
            id,
            json!({ "pong": true, "version": VERSION, "pid": std::process::id() }),
        ),
        "sidecar.search" => match req.params.map(serde_json::from_value::<SearchParams>) {
            Some(Ok(p)) => Response::ok(id, json!(search(&p.tools, &p.query, p.limit))),
            Some(Err(e)) => Response::fail(id, code::BAD_REQUEST, format!("invalid params: {e}")),
            None => Response::fail(id, code::BAD_REQUEST, "missing params"),
        },
        other => Response::fail(
            id,
            code::METHOD_NOT_FOUND,
            format!("unknown method {other:?}"),
        ),
    }
}

fn main() -> io::Result<()> {
    let stdin = io::stdin();
    let mut input = stdin.lock();
    let mut out = BufWriter::new(io::stdout().lock());

    while let Some(buf) = read_frame(&mut input)? {
        let resp = match serde_json::from_slice::<Request>(&buf) {
            Ok(req) => handle(req),
            Err(e) => Response::fail("", code::BAD_REQUEST, format!("invalid request: {e}")),
        };
        let bytes = serde_json::to_vec(&resp).unwrap_or_else(|_| {
            br#"{"id":"","ok":false,"error":{"code":"internal","message":"encode failed"}}"#
                .to_vec()
        });
        write_frame(&mut out, &bytes)?;
    }
    Ok(())
}
