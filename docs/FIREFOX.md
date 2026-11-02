# Firefox

- Build: `pnpm build:firefox` (MV2 por defecto en WXT) / `pnpm zip:firefox` (genera tambien el zip de fuentes que pide AMO).
- `SOURCE_CODE_REVIEW.md` ya trae las instrucciones de build para los revisores. Prueba el zip de fuentes: debe reconstruir una salida identica.
- Cuidado con `.env`: puede cambiar los hashes de los chunks; excluyelo del zip o usa `zip.includeSources`.
- Native messaging: `allowed_extensions` con el ID de `browser_specific_settings.gecko.id` (`toolbox@example.com` en `wxt.config.ts`; cambialo por el tuyo).
  Firefox lanza el host pasando la ruta del manifest y el ID del add-on; en Windows lo mete en un Job object (los subprocesos que deban sobrevivir necesitan `CREATE_BREAKAWAY_FROM_JOB`).
- APIs propias: ver `docs/API-COVERAGE.md` (`src/services/firefox/`). Los tipos de WXT salen de `@types/chrome`: para `sidebarAction`, `menus`, etc. instala `@wxt-dev/browser` + `@types/firefox-webext-browser`.
- No hay paginas `sandbox` ni `offscreen` en Firefox: ya estan excluidas con `manifest.exclude`.
- Fuente de verdad de la implementacion: https://searchfox.org/firefox-main/source/toolkit/components/extensions/ (la pagina de directorio no se pudo leer desde aqui; consultala a mano para comportamientos finos).
