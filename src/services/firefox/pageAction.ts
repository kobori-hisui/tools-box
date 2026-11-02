import { has } from '@/platform/capabilities';

/** browser.pageAction - Address-bar button. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'pageAction';

export const isAvailable = (): boolean => has(NAMESPACE);
