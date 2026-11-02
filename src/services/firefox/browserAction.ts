import { has } from '@/platform/capabilities';

/** browser.browserAction - Toolbar button (MV2 key). Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'browserAction';

export const isAvailable = (): boolean => has(NAMESPACE);
