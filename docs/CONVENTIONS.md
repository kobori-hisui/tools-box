# Convenciones

- Imports con `@/` (src) y `@@/` (raiz). Evita rutas relativas largas.
- Los entrypoints no contienen logica: importan de `tools/`, `core/`, `components/`.
- Una tool nunca importa de otra tool; lo compartido sube a `core/`, `services/` o `utils/`.
- Permisos amplios se piden como `optional_permissions` en runtime, por tool.
- Storage de cada tool con prefijo `area:<tool-id>:clave` y `version` + migraciones.
- Mensajes entre contextos: un unico `ProtocolMap` en `core/messaging/protocol.ts`.
- Todo texto visible via `public/_locales/*/messages.json`.
- Firefox: revisar `SOURCE_CODE_REVIEW.md` y probar el sources zip.
