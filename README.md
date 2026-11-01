# Toolbox (WXT + TypeScript + React)

Extension modular pensada para crecer a miles de herramientas sin tocar el nucleo.

## Comandos
```sh
pnpm install
pnpm dev            # Chrome   (pnpm dev:firefox / dev:edge)
pnpm build          # .output/chrome-mv3
pnpm zip:firefox    # zip + sources zip
pnpm tool:new developer base64 "Base64 Encoder"
pnpm tool:list
pnpm test && pnpm test:e2e
```

## Mapa del proyecto
| Carpeta | Rol |
|---|---|
| `src/entrypoints/` | Solo puntos de entrada **finos** (max. 1 nivel). Sin logica de negocio. |
| `src/tools/categories/<cat>/<tool>/` | **Una carpeta por herramienta**, autocontenida. |
| `src/tools/_core/` | Tipos, `defineTool`, registry (descubre tools con `import.meta.glob`). |
| `src/tools/_template/` | Plantilla que copia `pnpm tool:new`. |
| `src/core/` | Mensajeria, storage, settings, permisos, logger, flags, analytics. |
| `src/services/` | Wrappers de APIs del navegador y servicios externos (IA, API). |
| `src/platform/` | Diferencias entre navegadores (Chrome / Firefox / Safari). |
| `src/components`, `hooks`, `utils` | UI y helpers compartidos (dirs auto-importados por WXT). |
| `src/shared/` | Tipos, constantes y schemas compartidos. |
| `modules/` | Modulos WXT locales (build-time, se cargan solos). |
| `public/` | Archivos copiados tal cual (`_locales`, iconos, wasm, rules). |
| `src/_optional-entrypoints/` | newtab/history/bookmarks: mover a `entrypoints/` para activar. |

Ver `docs/` para arquitectura y convenciones.
