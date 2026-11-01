# Agregar una herramienta

```sh
pnpm tool:new <categoria> <tool-id> "Nombre"
```

Edita en la carpeta creada:
1. `meta.ts` – id, keywords, surfaces, `matches` (si corre en paginas), permisos opcionales.
2. `ui/Panel.tsx` – UI para popup/sidepanel.
3. `background.ts` / `content.ts` – solo si hacen falta.
4. `messages.ts`, `storage.ts`, `settings.ts` – estado propio (prefijo `local:<tool-id>:`).
5. `__tests__/` – tests con Vitest (usa `fakeBrowser`).

No hay que registrar nada: el registry la descubre sola.
