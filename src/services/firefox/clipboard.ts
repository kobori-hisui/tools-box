import { has } from '@/platform/capabilities';

/** browser.clipboard - Copy images to the clipboard. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'clipboard';

export const isAvailable = (): boolean => has(NAMESPACE);
