# protocol/

`fixtures/` son los golden files compartidos por TypeScript (`tests/unit`), Go (`native/go/pkg/protocol`, `internal/host`)
y Rust (`crates/protocol/tests`, `crates/core`). Ver `docs/PROTOCOL.md`.

`search.cases.json` define el comportamiento exacto de la busqueda; `src/core/search/search.ts` y `crates/core/src/search.rs` deben cumplirlo.
