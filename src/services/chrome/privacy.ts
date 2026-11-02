import { has } from '@/platform/capabilities';

/** chrome.privacy - Privacy-related settings. */
export const NAMESPACE = 'privacy';

/** False when the browser lacks the API or the permission was not granted. */
export const isAvailable = (): boolean => has(NAMESPACE);

// Put typed wrappers here. Tools must call these helpers, never the raw API,
// so browser differences stay in one place.
