# Arquitectura

## Idea central
Cada herramienta tiene **meta liviana** (`meta.ts`, eager) e **implementacion pesada** (`index.ts`, lazy).
El registry las descubre con `import.meta.glob`, asi el popup lista miles de tools sin cargar su codigo.

## Flujo
1. `popup / sidepanel / options` -> `ToolShell` -> `searchTools()` (solo metadata).
2. Al abrir una tool -> `loadToolImpl(id)` -> import dinamico de su `Panel`.
3. Content: `tools.content` (loader unico) -> carga `content.ts` de las tools cuyo `matches` coincide con la URL.
4. Background: `entrypoints/background/*` registra listeners (sincrono, dentro de `main`) y delega a `tool.onBackground`.
5. Mundo MAIN: `injectScript('/main-world.js')` desde el content script.

## Reglas de WXT que condicionan el diseno
- `entrypoints/` solo admite 0-1 nivel de profundidad: por eso las tools viven fuera, en `src/tools/`.
- El codigo de runtime de background/content va **dentro de `main`** (WXT importa el archivo en Node).
- Los content scripts de dev se registran dinamicamente (no aparecen en el manifest durante `wxt dev`).
- Assets de `public/` no son accesibles desde content scripts sin `web_accessible_resources`.
- UIs en paginas: `createShadowRootUi` (estilos aislados); usa px en vez de rem (ver FAQ de WXT).
