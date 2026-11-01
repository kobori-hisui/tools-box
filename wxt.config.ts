import { defineConfig } from 'wxt';

// Documentacion: https://wxt.dev/guide/essentials/config/manifest
export default defineConfig({
  srcDir: 'src', // codigo en src/, modules/ y public/ se quedan en la raiz
  modules: ['@wxt-dev/module-react'],

  manifest: ({ browser }) => ({
    name: '__MSG_extName__',
    description: '__MSG_extDescription__',
    default_locale: 'en',
    permissions: [
      'storage',
      'activeTab',
      'scripting',
      'contextMenus',
      'alarms',
      'downloads',
      'notifications',
      'clipboardWrite',
    ],
    // Permisos pedidos bajo demanda (por herramienta) en vez de al instalar
    optional_permissions: ['tabs', 'history', 'bookmarks', 'cookies', 'webNavigation'],
    optional_host_permissions: ['<all_urls>'],
    commands: {
      'open-command-palette': {
        suggested_key: { default: 'Ctrl+Shift+K', mac: 'Command+Shift+K' },
        description: 'Open command palette',
      },
    },
    // Script inyectado al mundo MAIN via injectScript() (ver entrypoints/main-world)
    web_accessible_resources: [
      { resources: ['main-world.js'], matches: ['<all_urls>'] },
    ],
    ...(browser === 'firefox'
      ? { browser_specific_settings: { gecko: { id: 'toolbox@example.com' } } }
      : {}),
  }),

  zip: {
    excludeSources: ['store/**', 'tests/**', 'docs/**'],
  },

  // Si WSL/Docker no detecta cambios: descomenta
  // watchOptions: { usePolling: true, interval: 1000 },
});
