/**
 * Runtime feature detection. Never branch on the browser name when you can
 * test for the API: Chrome, Firefox, Edge and Safari each lack different ones,
 * and an API can also be missing because its permission was not granted.
 */
function root(): unknown {
  const g = globalThis as { browser?: unknown; chrome?: unknown };
  return g.browser ?? g.chrome;
}

/** True when `path` (e.g. "sidePanel.open") exists on the extension API. */
export function has(path: string): boolean {
  let cur: unknown = root();
  for (const key of path.split('.')) {
    if (cur == null || typeof cur !== 'object') return false;
    cur = (cur as Record<string, unknown>)[key];
    if (cur == null) return false;
  }
  return true;
}

/** Features tools may depend on, mapped to the API that proves support. */
export const CAPABILITIES = {
  sidePanel: 'sidePanel.open', // Chrome / Edge
  sidebarAction: 'sidebarAction.open', // Firefox
  offscreen: 'offscreen.createDocument', // Chrome / Edge
  declarativeNetRequest: 'declarativeNetRequest.updateDynamicRules',
  userScripts: 'userScripts.register',
  nativeMessaging: 'runtime.connectNative',
  tabGroups: 'tabGroups.query',
  readingList: 'readingList.query',
  contextualIdentities: 'contextualIdentities.query', // Firefox containers
  tabCapture: 'tabCapture.capture',
  debugger: 'debugger.attach',
  webRequestBlocking: 'webRequest.onBeforeRequest',
  scripting: 'scripting.executeScript',
  storageSession: 'storage.session',
  alarms: 'alarms.create',
  notifications: 'notifications.create',
  downloads: 'downloads.download',
  omnibox: 'omnibox.onInputChanged',
  commands: 'commands.onCommand',
} as const;

export type Capability = keyof typeof CAPABILITIES;

export const supports = (c: Capability): boolean => has(CAPABILITIES[c]);

/** Snapshot of every capability (handy for a diagnostics screen / bug reports). */
export function snapshot(): Record<Capability, boolean> {
  return Object.fromEntries(
    (Object.keys(CAPABILITIES) as Capability[]).map((c) => [c, supports(c)]),
  ) as Record<Capability, boolean>;
}
