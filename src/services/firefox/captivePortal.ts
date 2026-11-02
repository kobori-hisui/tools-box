import { has } from '@/platform/capabilities';

/** browser.captivePortal - Captive portal state. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'captivePortal';

export const isAvailable = (): boolean => has(NAMESPACE);
