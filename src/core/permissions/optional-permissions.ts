import { browser } from 'wxt/browser';

/**
 * Request optional permissions. MUST be called from a user gesture (click).
 * `optional_permissions` / `optional_host_permissions` in wxt.config.ts must
 * list everything requested here.
 */
export async function ensurePermissions(opts: {
  permissions?: string[];
  origins?: string[];
}): Promise<boolean> {
  const req = opts as Parameters<typeof browser.permissions.request>[0];
  if (await browser.permissions.contains(req)) return true;
  return browser.permissions.request(req);
}

export async function dropPermissions(opts: {
  permissions?: string[];
  origins?: string[];
}): Promise<boolean> {
  return browser.permissions.remove(opts as Parameters<typeof browser.permissions.remove>[0]);
}
