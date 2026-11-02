import { has } from '@/platform/capabilities';

/** browser.browserSettings - Modify global browser settings. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'browserSettings';

export const isAvailable = (): boolean => has(NAMESPACE);
