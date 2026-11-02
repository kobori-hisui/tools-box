import { has } from '@/platform/capabilities';

/** browser.theme - Browser theme. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'theme';

export const isAvailable = (): boolean => has(NAMESPACE);
