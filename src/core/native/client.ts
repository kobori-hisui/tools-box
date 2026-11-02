import { browser } from 'wxt/browser';
import {
  isEvent,
  isResponse,
  NATIVE_HOST_NAME,
  NativeError,
  type NativeEvent,
  type NativeEvents,
  type NativeMethods,
} from './protocol';

/** The slice of runtime.Port the client needs (also what tests fake). */
export interface PortLike {
  postMessage(message: unknown): void;
  disconnect(): void;
  onMessage: { addListener(cb: (message: unknown) => void): void };
  onDisconnect: { addListener(cb: () => void): void };
}

export interface NativeClientOptions {
  hostName?: string;
  timeoutMs?: number;
  connect?: (hostName: string) => PortLike;
}

interface Pending {
  resolve(v: unknown): void;
  reject(e: unknown): void;
  timer: ReturnType<typeof setTimeout>;
}

/**
 * Talks to the Go native-messaging host over one long-lived port.
 * Connects lazily, matches responses by id, turns a dead host into
 * NativeError('unavailable') and reconnects on the next call.
 * Needs the `nativeMessaging` permission (optional - request it from a user gesture).
 */
export class NativeClient {
  private port: PortLike | undefined;
  private nextId = 0;
  private readonly pending = new Map<string, Pending>();
  private readonly listeners = new Set<(e: NativeEvent) => void>();
  private readonly hostName: string;
  private readonly timeoutMs: number;
  private readonly connectFn: (hostName: string) => PortLike;

  constructor(opts: NativeClientOptions = {}) {
    this.hostName = opts.hostName ?? NATIVE_HOST_NAME;
    this.timeoutMs = opts.timeoutMs ?? 30_000;
    this.connectFn =
      opts.connect ?? ((name) => browser.runtime.connectNative(name) as unknown as PortLike);
  }

  call<M extends keyof NativeMethods>(
    method: M,
    ...params: NativeMethods[M]['params'] extends undefined
      ? [params?: undefined]
      : [params: NativeMethods[M]['params']]
  ): Promise<NativeMethods[M]['result']>;
  call(method: string, params?: unknown): Promise<unknown>;
  call(method: string, params?: unknown): Promise<unknown> {
    let port: PortLike;
    try {
      port = this.ensurePort();
    } catch (e) {
      return Promise.reject(new NativeError('unavailable', `cannot start native host: ${String(e)}`));
    }
    const id = `c${++this.nextId}`;
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        this.pending.delete(id);
        reject(new NativeError('timeout', `${method} timed out after ${this.timeoutMs} ms`));
      }, this.timeoutMs);
      this.pending.set(id, { resolve, reject, timer });
      try {
        port.postMessage(params === undefined ? { id, method } : { id, method, params });
      } catch (e) {
        clearTimeout(timer);
        this.pending.delete(id);
        reject(new NativeError('unavailable', String(e)));
      }
    });
  }

  /** Subscribe to host-pushed events. Returns an unsubscribe function. */
  onEvent<E extends keyof NativeEvents>(
    name: E,
    cb: (data: NativeEvents[E]) => void,
  ): () => void {
    const listener = (e: NativeEvent) => {
      if (e.event === name) cb(e.data as NativeEvents[E]);
    };
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }

  close(): void {
    this.port?.disconnect();
    this.handleDisconnect();
  }

  private ensurePort(): PortLike {
    if (this.port) return this.port;
    const port = this.connectFn(this.hostName);
    port.onMessage.addListener((m) => this.handleMessage(m));
    port.onDisconnect.addListener(() => this.handleDisconnect());
    this.port = port;
    return port;
  }

  private handleMessage(m: unknown): void {
    if (isEvent(m)) {
      this.listeners.forEach((l) => l(m));
      return;
    }
    if (!isResponse(m)) return;
    const p = this.pending.get(m.id);
    if (!p) return;
    clearTimeout(p.timer);
    this.pending.delete(m.id);
    if (m.ok) p.resolve(m.result);
    else p.reject(new NativeError(m.error?.code ?? 'internal', m.error?.message ?? 'native error'));
  }

  private handleDisconnect(): void {
    this.port = undefined;
    for (const [id, p] of this.pending) {
      clearTimeout(p.timer);
      p.reject(new NativeError('unavailable', 'native host disconnected'));
      this.pending.delete(id);
    }
  }
}
