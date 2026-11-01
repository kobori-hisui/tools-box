import { defineAppConfig } from '#imports';

// Config de runtime (https://wxt.dev/guide/essentials/config/runtime)
declare module 'wxt/utils/define-app-config' {
  export interface WxtAppConfig {
    appName: string;
    featureFlags: Record<string, boolean>;
  }
}

export default defineAppConfig({
  appName: 'Toolbox',
  featureFlags: {
    commandPalette: true,
    aiTools: false,
  },
});
