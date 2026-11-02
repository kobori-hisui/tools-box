import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { NativeClient, NativeError, type PortLike } from '@/core/native';

class FakePort implements PortLike {
  sent: unknown[] = [];
  disconnected = false;
  private msg: Array<(m: unknown) => void> = [];
  private dis: Array<() => void> = [];
  onMessage = { addListener: (cb: (m: unknown) => void) => void this.msg.push(cb) };
  onDisconnect = { addListener: (cb: () => void) => void this.dis.push(cb) };
  postMessage(m: unknown) {
    this.sent.push(m);
  }
  disconnect() {
    this.disconnected = true;
  }
  emit(m: unknown) {
    this.msg.forEach((cb) => cb(m));
  }
  drop() {
    this.dis.forEach((cb) => cb());
  }
}

const fixture = (name: string) =>
  JSON.parse(readFileSync(resolve(__dirname, '../../protocol/fixtures', name), 'utf8'));

describe('NativeClient', () => {
  let ports: FakePort[];
  let client: NativeClient;
  beforeEach(() => {
    ports = [];
    client = new NativeClient({
      timeoutMs: 1000,
      connect: () => {
        const p = new FakePort();
        ports.push(p);
        return p;
      },
    });
  });
  afterEach(() => vi.useRealTimers());

  it('sends the same request shape as the golden fixture and resolves with the result', async () => {
    const p = client.call('ping');
    expect(ports[0]!.sent[0]).toEqual({ id: 'c1', method: 'ping' });
    const golden = fixture('response.ping.ok.json');
    ports[0]!.emit({ ...golden, id: 'c1' });
    await expect(p).resolves.toEqual(golden.result);
  });

  it('includes params when given', () => {
    void client.call('fs.watch', { path: '/tmp/inbox', recursive: true });
    expect(ports[0]!.sent[0]).toEqual({ id: 'c1', method: 'fs.watch', params: { path: '/tmp/inbox', recursive: true } });
  });

  it('rejects with the protocol error code', async () => {
    const p = client.call('fs.watch', { path: '/etc' });
    ports[0]!.emit({ ...fixture('response.error.json'), id: 'c1' });
    await expect(p).rejects.toMatchObject({ name: 'NativeError', code: 'forbidden' });
  });

  it('routes concurrent responses by id, out of order', async () => {
    const a = client.call('ping');
    const b = client.call('system.methods');
    ports[0]!.emit({ id: 'c2', ok: true, result: ['ping'] });
    ports[0]!.emit({ id: 'c1', ok: true, result: { pong: true } });
    await expect(b).resolves.toEqual(['ping']);
    await expect(a).resolves.toEqual({ pong: true });
    expect(ports).toHaveLength(1); // one shared port
  });

  it('delivers pushed events to subscribers only until unsubscribed', () => {
    const cb = vi.fn();
    void client.call('ping');
    const off = client.onEvent('fs.changed', cb);
    const ev = fixture('event.fs-changed.json');
    ports[0]!.emit(ev);
    expect(cb).toHaveBeenCalledWith(ev.data);
    off();
    ports[0]!.emit(ev);
    expect(cb).toHaveBeenCalledTimes(1);
  });

  it('times out', async () => {
    vi.useFakeTimers();
    const p = client.call('ping');
    const assertion = expect(p).rejects.toMatchObject({ code: 'timeout' });
    await vi.advanceTimersByTimeAsync(1001);
    await assertion;
  });

  it('fails pending calls on disconnect and reconnects on the next call', async () => {
    const p = client.call('ping');
    ports[0]!.drop();
    await expect(p).rejects.toMatchObject({ code: 'unavailable' });
    void client.call('ping');
    expect(ports).toHaveLength(2);
  });

  it('reports a failing connect as unavailable', async () => {
    const bad = new NativeClient({
      connect: () => {
        throw new Error('Specified native messaging host not found.');
      },
    });
    await expect(bad.call('ping')).rejects.toBeInstanceOf(NativeError);
    await expect(bad.call('ping')).rejects.toMatchObject({ code: 'unavailable' });
  });

  it('ignores responses for unknown ids', () => {
    void client.call('ping');
    expect(() => ports[0]!.emit({ id: 'zzz', ok: true })).not.toThrow();
  });
});
