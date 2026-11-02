//! Every golden file must parse into the Rust types and serialise back to the
//! same JSON value. Go and TypeScript run equivalent checks.

use serde_json::Value;
use std::path::Path;
use toolbox_protocol::{Event, Request, Response};

fn load(name: &str) -> String {
    let p = Path::new(env!("CARGO_MANIFEST_DIR"))
        .join("../../../../protocol/fixtures")
        .join(name);
    std::fs::read_to_string(&p).unwrap_or_else(|e| panic!("{}: {e}", p.display()))
}

fn round_trip<T>(name: &str)
where
    T: serde::de::DeserializeOwned + serde::Serialize,
{
    let raw = load(name);
    let original: Value = serde_json::from_str(&raw).unwrap();
    let typed: T = serde_json::from_str(&raw).unwrap_or_else(|e| panic!("{name}: {e}"));
    assert_eq!(
        serde_json::to_value(&typed).unwrap(),
        original,
        "{name} changed after round trip"
    );
}

#[test]
fn requests() {
    round_trip::<Request>("request.ping.json");
    round_trip::<Request>("request.fs-watch.json");
}

#[test]
fn responses() {
    round_trip::<Response>("response.ping.ok.json");
    round_trip::<Response>("response.error.json");
}

#[test]
fn events() {
    round_trip::<Event>("event.fs-changed.json");
}
