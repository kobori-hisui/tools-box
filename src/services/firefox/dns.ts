import { has } from '@/platform/capabilities';

/** browser.dns - Resolve domain names. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'dns';

export const isAvailable = (): boolean => has(NAMESPACE);
