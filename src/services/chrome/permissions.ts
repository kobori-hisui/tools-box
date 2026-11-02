import { has } from '@/platform/capabilities';

/** chrome.permissions - Optional permissions at runtime. */
export const NAMESPACE = 'permissions';

/** False when the browser lacks the API or the permission was not granted. */
export const isAvailable = (): boolean => has(NAMESPACE);

// Put typed wrappers here. Tools must call these helpers, never the raw API,
// so browser differences stay in one place.
