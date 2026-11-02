import { NativeError, type NativeEvent, type NativeMethods, type NativeResponse } from './protocol';

export interface DaemonOptions {
  baseUrl?: string; // default http://127.0.0.1:48731
  token: string; // sent as x-api-key
  fetch?: typeof fetch;
}

/**
 * Client for toolboxd (native/go/cmd/toolboxd): same protocol as the native
 * host, over loopback HTTP, for work that must outlive the browser.
 * Requires host permission for the daemon URL and its origin in config.json.
 */
export class DaemonClient {
  private readonly baseUrl: string;
  private readonly token: string;
  private readonly fetchFn: typeof fetch;
  private nextId = 0;

  constructor(opts: DaemonOptions) {
    this.baseUrl = (opts.baseUrl ?? 'http://127.0.0.1:48731').replace(/\/$/, '');
    this.token = opts.token;
    this.fetchFn = opts.fetch ?? ((...a) => fetch(...a));
  }

  async call<M extends keyof NativeMethods>(
    method: M,
    params?: NativeMethods[M]['params'],
  ): Promise<NativeMethods[M]['result']> {
    const body = JSON.stringify({ id: `d${++this.nextId}`, method, ...(params === undefined ? {} : { params }) });
    let res: Response;
    try {
      res = await this.fetchFn(`${this.baseUrl}/rpc`, {
        method: 'POST',
        headers: { 'content-type': 'application/json', 'x-api-key': this.token },
        body,
      });
    } catch (e) {
      throw new NativeError('unavailable', `daemon unreachable: ${String(e)}`);
    }
    if (res.status === 401) throw new NativeError('forbidden', 'daemon rejected the api key');
    if (res.status === 403) throw new NativeError('forbidden', 'daemon rejected this origin');
    const msg = (await res.json()) as NativeResponse<NativeMethods[M]['result']>;
    if (!msg.ok) throw new NativeError(msg.error?.code ?? 'internal', msg.error?.message ?? 'daemon error');
    return msg.result as NativeMethods[M]['result'];
  }

  /**
   * Stream server-sent events until `signal` aborts. Uses fetch (not
   * EventSource) because the api key must travel in a header.
   */
  async events(onEvent: (e: NativeEvent) => void, signal: AbortSignal): Promise<void> {
    const res = await this.fetchFn(`${this.baseUrl}/events`, { headers: { 'x-api-key': this.token }, signal });
    if (!res.ok || !res.body) throw new NativeError('unavailable', `events: HTTP ${res.status}`);
    const reader = res.body.pipeThrough(new TextDecoderStream()).getReader();
    let buf = '';
    for (;;) {
      const { value, done } = await reader.read();
      if (done) return;
      buf += value;
      let i: number;
      while ((i = buf.indexOf('\n\n')) >= 0) {
        const block = buf.slice(0, i);
        buf = buf.slice(i + 2);
        for (const line of block.split('\n')) {
          if (line.startsWith('data: ')) onEvent(JSON.parse(line.slice(6)) as NativeEvent);
        }
      }
    }
  }
}
