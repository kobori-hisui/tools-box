import { has } from '@/platform/capabilities';

/** browser.sidebarAction - Sidebar. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'sidebarAction';

export const isAvailable = (): boolean => has(NAMESPACE);
