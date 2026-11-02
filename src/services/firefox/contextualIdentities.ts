import { has } from '@/platform/capabilities';

/** browser.contextualIdentities - Containers. Firefox only: always guard with isAvailable(). */
export const NAMESPACE = 'contextualIdentities';

export const isAvailable = (): boolean => has(NAMESPACE);
