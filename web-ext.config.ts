import { defineWebExtConfig } from 'wxt';
import { resolve } from 'node:path';

// https://wxt.dev/guide/essentials/config/browser-startup
export default defineWebExtConfig({
  // Perfil persistente: conserva sesiones/logins entre reinicios del dev server
  chromiumProfile: resolve('.wxt/chrome-data'),
  keepProfileChanges: true,
  startUrls: ['https://example.com'],
});
