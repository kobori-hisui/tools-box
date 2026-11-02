import { describe, expect, it } from 'vitest';
import { WasmSearch, type ToolboxWasmExports } from '@/core/wasm/search-wasm';

/**
 * The real module is built from native/rust/crates/wasm. This test checks the
 * JS side of the ABI (allocation, 64-bit packing, freeing) against a fake
 * implementation of the same exports backed by a plain ArrayBuffer.
 */
function fakeExports(handler: (input: string) => string | null) {
  const memory = { buffer: new ArrayBuffer(4096) as ArrayBufferLike };
  let bump = 16;
  const freed: Array<[number, number]> = [];
  const x: ToolboxWasmExports = {
    memory,
    tb_alloc(len) {
      const p = bump;
      bump += Math.max(len, 1);
      return p;
    },
    tb_free(ptr, len) {
      freed.push([ptr, len]);
    },
    tb_search(ptr, len) {
      const input = new TextDecoder().decode(new Uint8Array(memory.buffer, ptr, len));
      const out = handler(input);
      if (out === null) return 0n;
      const bytes = new TextEncoder().encode(out);
      const outPtr = x.tb_alloc(bytes.length);
      new Uint8Array(memory.buffer, outPtr, bytes.length).set(bytes);
      return (BigInt(outPtr) << 32n) | BigInt(bytes.length);
    },
  };
  return { x, freed };
}

describe('WasmSearch', () => {
  it('sends JSON in, decodes the packed (ptr,len) result and frees both buffers', () => {
    const seen: string[] = [];
    const { x, freed } = fakeExports((input) => {
      seen.push(input);
      return '["a","b"]';
    });
    const ids = new WasmSearch(x).search([{ id: 'a', name: 'Alpha' }], 'al', 2);
    expect(ids).toEqual(['a', 'b']);
    expect(JSON.parse(seen[0]!)).toEqual({ tools: [{ id: 'a', name: 'Alpha' }], query: 'al', limit: 2 });
    expect(freed).toHaveLength(2); // input buffer + output buffer
  });

  it('throws when the module reports failure (0)', () => {
    const { x, freed } = fakeExports(() => null);
    expect(() => new WasmSearch(x).search([], 'x')).toThrow(/wasm search failed/);
    expect(freed).toHaveLength(1); // input still freed
  });
});
