import { has } from '@/platform/capabilities';

/** chrome.dom - Extension-only DOM APIs. */
export const NAMESPACE = 'dom';

/** False when the browser lacks the API or the permission was not granted. */
export const isAvailable = (): boolean => has(NAMESPACE);

// Put typed wrappers here. Tools must call these helpers, never the raw API,
// so browser differences stay in one place.
