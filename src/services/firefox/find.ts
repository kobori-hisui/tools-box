import { has } from '@/platform/capabilities';

/** browser.find - Find and highlight text. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'find';

export const isAvailable = (): boolean => has(NAMESPACE);
