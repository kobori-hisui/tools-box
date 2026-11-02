# Capas nativas: donde va cada cosa

```
 extension (TypeScript)                              tu maquina
┌────────────────────────────┐
│ popup / sidepanel / content│
│        │ sendMessage       │
│ background (service worker)│── connectNative ──► toolbox-host (Go)  ── stdin/stdout ──► toolbox-sidecar (Rust)
│   core/native/client.ts    │   (lo lanza el navegador)   │ fs.watch, notify.show           │ busqueda, calculo pesado,
│   core/native/daemon.ts    │── fetch 127.0.0.1 ─► toolboxd (Go, siempre encendido)          │ drivers (serial/HID/USB)
│                            │   x-api-key + origen        │ mismo router, SSE de eventos
│ core/wasm (Rust -> .wasm)  │  corre DENTRO del navegador, sin instalar nada
└────────────────────────────┘
```

## Que lenguaje para que

| Necesidad | Donde | Por que |
|---|---|---|
| Calculo puro (busqueda, parseo, hashing, compresion) | `native/rust/crates/core` -> **WASM** (`crates/wasm`) | Va dentro de la extension: sin instalador, sin permisos, funciona en Chrome/Firefox/Safari. |
| Lo mismo pero enorme o que necesita hilos/disco | **sidecar Rust** (`crates/sidecar`) | Sin limites de WASM; lo lanza el host Go. |
| Notificaciones del SO, vigilar carpetas, procesos, red local | **host Go** (`native/go`) | Binario unico por OS, compila cruzado facil, ya conoces el stack. |
| Que siga vivo con el navegador cerrado (watchers, tareas programadas) | **daemon Go** `toolboxd` | HTTP loopback + SSE; mismo router que el host. |
| Hardware de bajo nivel (serial, HID, USB, MMIO) | `native/rust/crates/drivers` (rapido/unsafe) o `native/go/internal/drivers/*` (cgo) | Aislado: sus dependencias nunca llegan al WASM. |

Regla: **lo que pueda ir en WASM o TypeScript no pide instalar nada**. El host nativo es opcional (`native:available`) y cada tool declara en `meta.ts` si lo necesita (`native: { methods: [...] }`).

## Donde se agrega cada cosa

| Quieres... | Haz esto |
|---|---|
| Un metodo nuevo del host (ej. `fs.list`) | `native/go/internal/app/app.go` -> `r.Register("fs.list", ...)`; tipos en `src/core/native/protocol.ts` (`NativeMethods`); fixture en `protocol/fixtures/`. |
| Un metodo de Rust | Maneja `sidecar.<nombre>` en `crates/sidecar/src/main.rs`; llega solo, el host reenvia el prefijo `sidecar.`. |
| Logica reutilizable en WASM y sidecar | `crates/core`; si es para JS exporta una funcion C-ABI en `crates/wasm`. |
| Un driver | `crates/drivers` (Rust) o `internal/drivers/<familia>` (Go) + metodo `drivers.<familia>.*`. |
| Notificaciones en Windows | `internal/notify/notify.go`: hoy devuelve `unavailable` en Windows (los toast necesitan AUMID/WinRT). Implementalo ahi. |
| Notificaciones por eventos del kernel (inotify/FSEvents) | Implementa `watcher.Backend` con `github.com/fsnotify/fsnotify` y pasalo a `watcher.NewManager`. Hoy es polling stdlib (correcto en todos lados, menos eficiente). |

## Seguridad (ya aplicada)

- El host **deniega todo** por defecto: `fs.*` solo funciona dentro de `allowedRoots` (resuelve symlinks, bloquea `..`).
- Chrome solo deja que los IDs de `allowed_origins` hablen con el host; Firefox usa `allowed_extensions`.
- El daemon solo escucha en `127.0.0.1`, exige `x-api-key` (>= 16 chars, comparacion en tiempo constante) y, si viene `Origin`, debe estar en la lista.
- Respuestas > 1 MB se convierten en error `too_large` (el navegador cortaria la conexion).
- Nada se ejecuta via shell: los argumentos van como argv.

## Instalar el host

```sh
pnpm native:build                       # native/dist/<os>-<arch>/
pnpm native:install -- --chrome-id <ID> --firefox-id toolbox@example.com
pnpm native:install -- --print          # ver los manifests sin escribir nada
pnpm native:install -- --uninstall
```
Copia `native/config.example.json` a `<config del usuario>/toolbox/config.json` (Linux `~/.config/toolbox/`, macOS `~/Library/Application Support/toolbox/`, Windows `%AppData%\toolbox\`).

Los directorios/claves de registro de Firefox en Windows estan confirmados por MDN. Los de Chrome/Edge/Chromium
siguen el layout documentado de Chrome pero **no pude abrir esa pagina** durante el diseno: verifica con `--print`.
En Chrome el ID de la extension solo es estable si fijas `key` en el manifest o tras publicar.
