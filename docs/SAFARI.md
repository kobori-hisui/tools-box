# Safari

Fuente: documentacion de Apple (articulo "Converting a Safari app extension to a Safari web extension") y la guia
de publicacion de WXT. No pude leer la pagina principal de Apple ni la de compatibilidad (el fetch fue bloqueado),
asi que la matriz de APIs soportadas **no esta verificada aqui**.

## Lo confirmado
- Safari 14+ (macOS) y 15+ (iOS) soportan web extensions con las mismas APIs JavaScript que otros navegadores.
- Safari usa `browser.*` con Promises desde siempre (igual que Firefox).
- Los IDs de recursos accesibles cambian por sesion en Safari; no los hardcodees.
- WXT no automatiza la publicacion en Safari: requiere una app contenedora nativa de macOS/iOS.
- Comando documentado por WXT: `pnpm wxt build -b safari` y luego `xcrun safari-web-extension-packager .output/safari-mv2` (pasa el directorio de salida, no el codigo fuente).

## Flujo
```sh
pnpm safari:convert            # build -b safari + xcrun safari-web-extension-packager (solo macOS)
```
Abre el proyecto de Xcode generado, firma con tu cuenta de desarrollador y publica la app.

## Native messaging en Safari
Safari **no ejecuta binarios arbitrarios** como Chrome/Firefox. Los mensajes (`runtime.sendNativeMessage`) llegan a un handler Swift
dentro de la app contenedora (`native/safari/SafariWebExtensionHandler.swift`, plantilla sin compilar). Para reutilizar Go/Rust:
la app contenedora lanza el host/daemon, o la extension habla con `toolboxd` en `127.0.0.1` (si el permiso de host lo permite).

## Regla de oro
Ramifica por capacidad (`supports('...')`), no por navegador. Marca en `meta.ts` `browsers: [...]` solo cuando una tool realmente no puede funcionar.
