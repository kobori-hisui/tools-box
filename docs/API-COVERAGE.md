# Cobertura de APIs de extensiones

Generado a partir de las referencias oficiales leidas al disenar el proyecto:
- Chrome: https://developer.chrome.com/docs/extensions/reference/api
- Firefox/MDN: https://developer.mozilla.org/en-US/Add-ons/WebExtensions/API/
- Safari: https://developer.apple.com/documentation/safariservices/safari-web-extensions

Cada namespace tiene un wrapper en `src/services/` con `isAvailable()` (deteccion en runtime).
Los tools usan esos wrappers y nunca la API cruda: asi las diferencias entre navegadores viven en un solo sitio.

## Chrome / Chromium (59 namespaces con wrapper)

| Namespace | Que hace | Wrapper |
|---|---|---|
| `accessibilityFeatures` | Manage Chrome's accessibility features - needs accessibilityFeatures.read/modify | `src/services/chrome/accessibilityFeatures.ts` |
| `action` | Toolbar icon (MV3) | `src/services/chrome/action.ts` |
| `alarms` | Schedule code | `src/services/chrome/alarms.ts` |
| `bookmarks` | Create and organize bookmarks | `src/services/chrome/bookmarks.ts` |
| `browsingData` | Remove browsing data | `src/services/chrome/browsingData.ts` |
| `commands` | Keyboard shortcuts | `src/services/chrome/commands.ts` |
| `contentSettings` | Per-site content settings | `src/services/chrome/contentSettings.ts` |
| `contextMenus` | Context menu items | `src/services/chrome/contextMenus.ts` |
| `cookies` | Query and modify cookies | `src/services/chrome/cookies.ts` |
| `debugger` | Remote debugging protocol access | `src/services/chrome/debugger.ts` |
| `declarativeContent` | Act on page content without reading it | `src/services/chrome/declarativeContent.ts` |
| `declarativeNetRequest` | Block/modify requests with rules | `src/services/chrome/declarativeNetRequest.ts` |
| `desktopCapture` | Capture screen, windows or tabs | `src/services/chrome/desktopCapture.ts` |
| `devtools` | DevTools panels, network, inspectedWindow | `src/services/chrome/devtools.ts` |
| `dns` | DNS resolution - dev channel only | `src/services/chrome/dns.ts` |
| `dom` | Extension-only DOM APIs | `src/services/chrome/dom.ts` |
| `downloads` | Start and manage downloads | `src/services/chrome/downloads.ts` |
| `extension` | Extension utilities | `src/services/chrome/extension.ts` |
| `fontSettings` | Manage font settings | `src/services/chrome/fontSettings.ts` |
| `gcm` | Firebase Cloud Messaging | `src/services/chrome/gcm.ts` |
| `history` | Browser history | `src/services/chrome/history.ts` |
| `i18n` | Internationalization | `src/services/chrome/i18n.ts` |
| `identity` | OAuth2 tokens | `src/services/chrome/identity.ts` |
| `idle` | Detect idle state | `src/services/chrome/idle.ts` |
| `instanceID` | Instance ID service | `src/services/chrome/instanceID.ts` |
| `management` | Manage installed extensions | `src/services/chrome/management.ts` |
| `notifications` | Rich notifications | `src/services/chrome/notifications.ts` |
| `offscreen` | Offscreen documents (MV3) | `src/services/chrome/offscreen.ts` |
| `omnibox` | Address-bar keyword | `src/services/chrome/omnibox.ts` |
| `pageCapture` | Save a tab as MHTML | `src/services/chrome/pageCapture.ts` |
| `permissions` | Optional permissions at runtime | `src/services/chrome/permissions.ts` |
| `power` | Override power management | `src/services/chrome/power.ts` |
| `printerProvider` | Expose printers | `src/services/chrome/printerProvider.ts` |
| `privacy` | Privacy-related settings | `src/services/chrome/privacy.ts` |
| `processes` | Browser processes - dev channel only | `src/services/chrome/processes.ts` |
| `proxy` | Proxy settings | `src/services/chrome/proxy.ts` |
| `publicSuffix` | Public Suffix List queries - new (Chrome 153+) | `src/services/chrome/publicSuffix.ts` |
| `readingList` | Reading List (MV3) | `src/services/chrome/readingList.ts` |
| `runtime` | Service worker, manifest, lifecycle, native messaging | `src/services/chrome/runtime.ts` |
| `scripting` | Inject scripts and CSS (MV3) | `src/services/chrome/scripting.ts` |
| `search` | Search with the default provider | `src/services/chrome/search.ts` |
| `sessions` | Restore tabs and windows | `src/services/chrome/sessions.ts` |
| `sidePanel` | Browser side panel (MV3) | `src/services/chrome/sidePanel.ts` |
| `storage` | Store and watch data | `src/services/chrome/storage.ts` |
| `system.cpu` | CPU metadata | `src/services/chrome/system-cpu.ts` |
| `system.display` | Display metadata | `src/services/chrome/system-display.ts` |
| `system.memory` | Memory information | `src/services/chrome/system-memory.ts` |
| `system.storage` | Storage devices | `src/services/chrome/system-storage.ts` |
| `tabCapture` | Tab media streams | `src/services/chrome/tabCapture.ts` |
| `tabGroups` | Tab groups (MV3) | `src/services/chrome/tabGroups.ts` |
| `tabs` | Create and modify tabs | `src/services/chrome/tabs.ts` |
| `topSites` | Most-visited sites | `src/services/chrome/topSites.ts` |
| `tts` | Text to speech | `src/services/chrome/tts.ts` |
| `ttsEngine` | Implement a TTS engine | `src/services/chrome/ttsEngine.ts` |
| `userScripts` | User scripts world (MV3) | `src/services/chrome/userScripts.ts` |
| `webAuthenticationProxy` | Intercept WebAuthn (MV3) | `src/services/chrome/webAuthenticationProxy.ts` |
| `webNavigation` | Navigation events | `src/services/chrome/webNavigation.ts` |
| `webRequest` | Observe/modify in-flight requests | `src/services/chrome/webRequest.ts` |
| `windows` | Create and modify windows | `src/services/chrome/windows.ts` |

## Solo Firefox (12 namespaces con wrapper)

| Namespace | Que hace | Wrapper |
|---|---|---|
| `browserAction` | Toolbar button (MV2 key) | `src/services/firefox/browserAction.ts` |
| `browserSettings` | Modify global browser settings | `src/services/firefox/browserSettings.ts` |
| `captivePortal` | Captive portal state | `src/services/firefox/captivePortal.ts` |
| `clipboard` | Copy images to the clipboard | `src/services/firefox/clipboard.ts` |
| `contextualIdentities` | Containers | `src/services/firefox/contextualIdentities.ts` |
| `dns` | Resolve domain names | `src/services/firefox/dns.ts` |
| `find` | Find and highlight text | `src/services/firefox/find.ts` |
| `menus` | Menu system (superset of contextMenus) | `src/services/firefox/menus.ts` |
| `pageAction` | Address-bar button | `src/services/firefox/pageAction.ts` |
| `pkcs11` | PKCS #11 modules | `src/services/firefox/pkcs11.ts` |
| `sidebarAction` | Sidebar | `src/services/firefox/sidebarAction.ts` |
| `theme` | Browser theme | `src/services/firefox/theme.ts` |

Notas de Firefox:
- Usa `browser.*` con Promises; tambien acepta `chrome.*` con callbacks.
- `menus` es el superconjunto de `contextMenus`; `sidebarAction` equivale a `sidePanel` de Chrome.
- Firefox y Safari mantienen background page (MV2) por defecto en WXT; Chromium usa service worker (MV3).
- Los tipos de WXT salen de `@types/chrome`; para los namespaces solo-Firefox anade `@wxt-dev/browser` + `@types/firefox-webext-browser` (ver WXT: Extension APIs).

## Fuera de alcance a proposito

APIs exclusivas de ChromeOS / empresa (requieren politica o solo existen en ChromeOS): `audio`, `certificateProvider`, `documentScan`, `enterprise.deviceAttributes`, `enterprise.hardwarePlatform`, `enterprise.login`, `enterprise.networkingAttributes`, `enterprise.platformKeys`, `fileBrowserHandler`, `fileSystemProvider`, `input.ime`, `loginState`, `platformKeys`, `printing`, `printingMetrics`, `systemLog`, `vpnProvider`, `wallpaper`.
Si algun dia haces una variante para ChromeOS/empresa, crea `src/services/chromeos/`.

## Safari

Safari 14+ (macOS) y 15+ (iOS) soportan las mismas APIs JavaScript de web extensions que otros navegadores,
pero no todas. Comprueba la matriz oficial antes de depender de una API (`docs/SAFARI.md`).
`supports()` en `src/platform/capabilities.ts` es la forma segura de ramificar.

## Como saber que APIs tiene un navegador concreto

`snapshot()` (en `src/platform/capabilities.ts`) devuelve un mapa capacidad -> boolean. Muestralo en una pantalla de
diagnostico o ejecutalo desde la consola del service worker del navegador que quieras probar. En esa misma consola,
`Object.keys(chrome).sort()` lista todos los namespaces realmente presentes.
