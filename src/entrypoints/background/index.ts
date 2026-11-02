import { registerAlarms } from './alarms';
import { registerCommands } from './commands';
import { registerContextMenus } from './context-menus';
import { registerLifecycle } from './lifecycle';
import { registerMessaging } from './messaging';
import { registerNative } from './native';

export default defineBackground({
  type: 'module',
  main() {
    // TODO: todos los listeners deben registrarse aqui (sincrono, no async)
    registerLifecycle();
    registerMessaging();
    registerNative();
    registerContextMenus();
    registerAlarms();
    registerCommands();
  },
});
