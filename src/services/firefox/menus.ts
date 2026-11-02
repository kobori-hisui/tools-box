import { has } from '@/platform/capabilities';

/** browser.menus - Menu system (superset of contextMenus). Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'menus';

export const isAvailable = (): boolean => has(NAMESPACE);
