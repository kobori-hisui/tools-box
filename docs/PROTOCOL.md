# Protocolo extension <-> host / daemon / sidecar

Una sola definicion para los tres lenguajes. La fuente de verdad son los **golden files** de `protocol/fixtures/`:
TypeScript, Go y Rust los parsean en sus tests; si un tipo se desvia, el test falla.

## Mensajes (JSON UTF-8)

```jsonc
// peticion                { "id": "1", "method": "ping", "params": {...} }
// respuesta ok            { "id": "1", "ok": true,  "result": {...} }
// respuesta error         { "id": "1", "ok": false, "error": { "code": "forbidden", "message": "..." } }
// evento (sin id)         { "event": "fs.changed", "data": {...} }
```
Codigos: `bad_request`, `method_not_found`, `forbidden`, `too_large`, `timeout`, `internal`, `unavailable`.

## Transporte

| Canal | Framing |
|---|---|
| Native messaging (navegador <-> host) y host <-> sidecar | 4 bytes de longitud sin signo en **orden nativo** + JSON. Host -> navegador maximo 1 MB; navegador -> host mucho mayor (el host acepta hasta 64 MB). stdout solo lleva frames; logs a stderr. |
| Daemon | `POST /rpc` (JSON), `GET /events` (SSE), `GET /health`. Header `x-api-key`. |

Metodos actuales: `ping`, `system.methods`, `fs.watch`, `fs.unwatch`, `notify.show`, `sidecar.ping`, `sidecar.search`.

## Agregar un metodo

1. Fixture(s) en `protocol/fixtures/`.
2. Handler en Go (`internal/app`) o Rust (`crates/sidecar`).
3. Tipo en `src/core/native/protocol.ts`.
4. Test que use el fixture en el lenguaje donde vive el handler.
