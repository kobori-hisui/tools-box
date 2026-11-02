# Changelog

Formato basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.1.0/) y [SemVer](https://semver.org/lang/es/).

## [0.1.0] - 2026-11-02

Primera version: esqueleto modular de la extension (WXT + TypeScript + React) preparado para cientos o miles de herramientas.

### Added

- **Estructura WXT**: `srcDir: 'src'` con entrypoints `background`, `popup`, `options`, `sidepanel`, `devtools`, `devtools-panel`, `onboarding`, `offscreen`, `sandbox`, `overlay.content`, `tools.content` y `main-world`.
- **Sistema de herramientas**: una carpeta por tool en `src/tools/categories/<categoria>/<tool>/` con `meta.ts` ligero (carga inmediata) e `index.ts` (carga diferida); registro en `src/tools/_core/registry.ts`, plantilla `_template` y comando `pnpm tool:new`.
- **Tools de ejemplo**: `json-formatter`, `word-counter` y `page-media`.
- **Busqueda** compartida (`src/core/search`) con puntuacion, y version WASM (`src/core/wasm`).
- **Cobertura de APIs**: 59 wrappers de Chrome (`src/services/chrome`) y 12 de Firefox (`src/services/firefox`), con `isAvailable()` por namespace; detector de capacidades en `src/platform/capabilities.ts`. Detalle en `docs/API-COVERAGE.md`.
- **Capa nativa Go** (`native/go`, solo libreria estandar): host de native messaging `toolbox-host`, daemon `toolboxd` (127.0.0.1:48731, API key, allowlist de origenes, SSE) e instalador.
- **Capa nativa Rust** (`native/rust`): workspace con `protocol`, `core`, `wasm` (ABI C), `sidecar` y `drivers`.
- **Protocolo compartido** (`protocol/fixtures`): golden files usados por TypeScript, Go y Rust, incluida la paridad del algoritmo de busqueda.
- **Soporte Firefox**: build MV2, `data_collection_permissions`, instrucciones de revision de codigo fuente (`SOURCE_CODE_REVIEW.md`).
- **Soporte Safari** (plantilla): `pnpm safari:convert` y `native/safari/SafariWebExtensionHandler.swift`.
- **Onboarding** (`/onboarding.html`): catalogo de 118 ideas de herramientas en 15 categorias, solo nombres, con buscador y chips de las APIs de extension que usaria cada una (marcadas C = Chrome, F = Firefox).
- **Scripts**: `native:build`, `native:test`, `native:install`, `safari:convert`, `check`.
- **Calidad**: ESLint, Vitest (39 tests), Playwright (fixtures E2E), CI de nativos (`.github/workflows/native.yml`).
- **Documentacion** en `docs/`: arquitectura, convenciones, como agregar una tool, permisos, publicacion, nativos, protocolo, Firefox, Safari e ideas.

### Changed

- `@vitejs/plugin-react` fijado a `^5.0.0` en `pnpm-workspace.yaml` por compatibilidad con Vite 7 (pnpm 10+ lee los overrides de ese archivo, no de `package.json`).

### Known limitations

- La matriz de APIs soportadas en **Safari no esta verificada**; el catalogo de onboarding no afirma soporte en Safari.
- `SafariWebExtensionHandler.swift` es una plantilla sin compilar ni probar.
- El catalogo de onboarding son solo nombres; no hay implementacion de esas herramientas.
- La cobertura de Rust a WASM no se compilo para `wasm32` en el entorno de desarrollo (target no disponible); el ABI esta probado nativamente.

[0.1.0]: ./CHANGELOG.md
