import { supports } from '@/platform/capabilities';

interface OffscreenApi {
  createDocument(opts: { url: string; reasons: string[]; justification: string }): Promise<void>;
  closeDocument(): Promise<void>;
  hasDocument?(): Promise<boolean>;
}

const api = (): OffscreenApi | undefined =>
  (globalThis as { chrome?: { offscreen?: OffscreenApi } }).chrome?.offscreen;

let creating: Promise<void> | undefined;

/**
 * Chrome/Edge only. Keeps ONE offscreen document alive for DOM work a service
 * worker cannot do (DOMParser, Canvas, audio, clipboard, Web Workers, WASM).
 * Firefox background pages already have a DOM, so this is a no-op there.
 */
export async function ensureOffscreen(
  reasons: string[] = ['WORKERS'],
  justification = 'Run heavy tool computations off the service worker',
): Promise<void> {
  if (!supports('offscreen')) return;
  const off = api()!;
  if (off.hasDocument && (await off.hasDocument())) return;
  creating ??= off
    .createDocument({ url: 'offscreen.html', reasons, justification })
    .finally(() => (creating = undefined));
  await creating;
}

export async function closeOffscreen(): Promise<void> {
  if (supports('offscreen')) await api()!.closeDocument();
}
