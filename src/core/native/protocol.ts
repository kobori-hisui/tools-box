/** Wire types shared with native/go and native/rust. Fixtures: /protocol/fixtures. */

export const NATIVE_HOST_NAME = 'com.toolbox.host';

export interface NativeRequest {
  id: string;
  method: string;
  params?: unknown;
}

export interface NativeErrorBody {
  code: NativeErrorCode;
  message: string;
}

export interface NativeResponse<T = unknown> {
  id: string;
  ok: boolean;
  result?: T;
  error?: NativeErrorBody;
}

export interface NativeEvent<T = unknown> {
  event: string;
  data?: T;
}

export type NativeErrorCode =
  | 'bad_request'
  | 'method_not_found'
  | 'forbidden'
  | 'too_large'
  | 'timeout'
  | 'internal'
  | 'unavailable';

/** Methods implemented by the Go host / daemon (see native/go/internal/app). */
export interface NativeMethods {
  ping: { params: undefined; result: { pong: true; version: string; os: string; arch: string } };
  'system.methods': { params: undefined; result: string[] };
  'fs.watch': { params: { path: string; recursive?: boolean }; result: { watchId: string } };
  'fs.unwatch': { params: { watchId: string }; result: { removed: boolean } };
  'notify.show': { params: { title: string; body?: string }; result: { shown: boolean } };
  'sidecar.ping': { params: undefined; result: { pong: true; version: string; pid: number } };
  'sidecar.search': {
    params: { tools: unknown[]; query: string; limit?: number };
    result: string[];
  };
}

export interface NativeEvents {
  'fs.changed': { watchId: string; path: string; op: 'create' | 'write' | 'remove' };
}

export class NativeError extends Error {
  constructor(
    public readonly code: NativeErrorCode,
    message: string,
  ) {
    super(message);
    this.name = 'NativeError';
  }
}

export const isEvent = (m: unknown): m is NativeEvent =>
  typeof m === 'object' && m !== null && 'event' in m && !('id' in m);

export const isResponse = (m: unknown): m is NativeResponse =>
  typeof m === 'object' && m !== null && 'id' in m && 'ok' in m;
