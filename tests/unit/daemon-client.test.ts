import { describe, expect, it, vi } from 'vitest';
import { DaemonClient } from '@/core/native';

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), { status, headers: { 'content-type': 'application/json' } });

describe('DaemonClient', () => {
  it('posts the request with the api key and returns the result', async () => {
    const f = vi.fn().mockResolvedValue(json({ id: 'd1', ok: true, result: { pong: true } }));
    const c = new DaemonClient({ token: 'secret-secret-secret', fetch: f });
    await expect(c.call('ping')).resolves.toEqual({ pong: true });
    const [url, init] = f.mock.calls[0]!;
    expect(url).toBe('http://127.0.0.1:48731/rpc');
    expect(init.headers['x-api-key']).toBe('secret-secret-secret');
    expect(JSON.parse(init.body)).toEqual({ id: 'd1', method: 'ping' });
  });

  it('maps HTTP 401 and protocol errors', async () => {
    const unauth = new DaemonClient({ token: 'x', fetch: vi.fn().mockResolvedValue(new Response('', { status: 401 })) });
    await expect(unauth.call('ping')).rejects.toMatchObject({ code: 'forbidden' });

    const err = new DaemonClient({
      token: 'x',
      fetch: vi.fn().mockResolvedValue(json({ id: 'd1', ok: false, error: { code: 'method_not_found', message: 'nope' } })),
    });
    await expect(err.call('ping')).rejects.toMatchObject({ code: 'method_not_found', message: 'nope' });
  });

  it('reports an unreachable daemon as unavailable', async () => {
    const c = new DaemonClient({ token: 'x', fetch: vi.fn().mockRejectedValue(new TypeError('Failed to fetch')) });
    await expect(c.call('ping')).rejects.toMatchObject({ code: 'unavailable' });
  });

  it('parses server-sent events split across chunks', async () => {
    const enc = new TextEncoder();
    const body = new ReadableStream<Uint8Array>({
      start(ctrl) {
        ctrl.enqueue(enc.encode(': connected\n\ndata: {"event":"fs.chan'));
        ctrl.enqueue(enc.encode('ged","data":{"path":"/a"}}\n\n'));
        ctrl.close();
      },
    });
    const c = new DaemonClient({ token: 'x', fetch: vi.fn().mockResolvedValue(new Response(body)) });
    const got: unknown[] = [];
    await c.events((e) => got.push(e), new AbortController().signal);
    expect(got).toEqual([{ event: 'fs.changed', data: { path: '/a' } }]);
  });
});
