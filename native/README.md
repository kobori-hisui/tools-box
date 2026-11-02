# native/

```
go/      host (native messaging), daemon (HTTP loopback), installer      -> cmd/ + internal/
rust/    workspace: protocol, core (logica pura), wasm (C-ABI), sidecar, drivers
manifests/  ejemplos de manifests del host (el instalador los genera)
safari/  plantilla Swift (macOS)
config.example.json
```
Ver `docs/NATIVE.md`. Comandos: `pnpm native:build`, `pnpm native:test`, `pnpm native:install`.

Go solo usa la libreria estandar a proposito (sin dependencias que descargar). Rust usa `serde` y `serde_json`.
