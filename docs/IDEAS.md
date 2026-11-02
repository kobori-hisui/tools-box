# Ideas para no dejar huecos

Cada fila dice **donde vive** y que ya existe en el esqueleto. Lo que dice "pendiente" es una idea, no codigo.

## Categorias de herramientas
| Categoria | Ejemplos | Capas |
|---|---|---|
| developer | JSON/YAML/TOML, base64, JWT, regex, diff, cron, UUID, timestamps | TS (+ WASM para volumen) |
| text | contadores, case, slug, lorem, markdown | TS |
| media | compresor de imagen, conversor, recorte, capturas (`tabCapture`, `desktopCapture`, `pageCapture`) | offscreen + WASM |
| downloaders | `downloads` + HLS/segmentos, colas con reintentos | background + host Go (+ ffmpeg/yt-dlp por `internal/process`) |
| security | hash, generador de contrasenas, auditor de cookies, detector de trackers | TS + `cookies`, `webRequest` |
| network | DNS, headers, redirecciones, proxy (`proxy`), latencia | `dns`/`proxy` + daemon |
| seo | auditor on-page, meta, schema | content script |
| design | eyedropper, contraste, tipografias, paletas | content script + overlay |
| data | CSV/JSON/SQL en navegador, graficos | WASM / sidecar |
| ai | resumen, traduccion, OCR local | `src/services/ai/*` + sidecar para modelos locales |
| automation | macros, reglas por sitio, atajos (`commands`), scheduler (`alarms`) | background + daemon para lo programado |
| accessibility | lector, contraste, foco | content script |
| converters | unidades, divisas, formatos, archivos | TS / sidecar |
| social | descargas de medios por sitio | content + host |

## Infraestructura que conviene agregar cuando la necesites (pendiente)
- **Command palette global** (ya hay overlay + comando `Ctrl+Shift+K`): conectarlo a `searchTools()`.
- **Pantalla de diagnostico**: `snapshot()` de capacidades + `native:available` + versiones, para soporte y bug reports.
- **Permisos por tool**: `meta.permissions` -> `ensurePermissions` al abrir la tool; revocar al desactivar.
- **Activar/desactivar tools** y favoritos (`core/storage/items.ts` ya tiene los items).
- **Modulo WXT que genere el catalogo** (`modules/`): lista JSON de metas en build para no cargar todos los `meta.ts` en runtime cuando pasen de ~1000.
- **Busqueda en WASM/worker** cuando el catalogo sea enorme (`core/wasm`, `workers/search.worker.ts`).
- **Telemetria opt-in**, flags remotos, actualizador de reglas (`public/rules/` + `declarativeNetRequest`).
- **E2E por tool** con Playwright (`tests/e2e/fixtures.ts` ya carga la extension).
- **Firma / reproducibilidad**: `zip` de Firefox + CI de release (`.github/workflows/release.yml`).
- **Windows**: notificaciones toast reales y servicio/arranque automatico de `toolboxd`.
- **Monorepo**: cuando haya varias extensiones, mueve `native/` y `protocol/` a paquetes compartidos.
