import type { Searchable } from '@/core/search/search';

/**
 * Loader for native/rust/crates/wasm (plain C ABI, no wasm-bindgen).
 * Build + copy: `pnpm native:build:wasm` -> public/wasm/toolbox.wasm
 *
 * Protocol: tb_alloc(len) -> ptr; write UTF-8 JSON; tb_search(ptr, len) ->
 * u64 packed as (outPtr << 32) | outLen (0 = error); read JSON; tb_free both.
 */
export interface ToolboxWasmExports {
  memory: { buffer: ArrayBufferLike };
  tb_alloc(len: number): number;
  tb_free(ptr: number, len: number): void;
  tb_search(ptr: number, len: number): bigint;
}

export class WasmSearch {
  private readonly enc = new TextEncoder();
  private readonly dec = new TextDecoder();

  constructor(private readonly x: ToolboxWasmExports) {}

  static async load(url: string): Promise<WasmSearch> {
    const bytes = await (await fetch(url)).arrayBuffer();
    const { instance } = await WebAssembly.instantiate(bytes, {});
    return new WasmSearch(instance.exports as unknown as ToolboxWasmExports);
  }

  search(tools: readonly Searchable[], query: string, limit = 0): string[] {
    const input = this.enc.encode(JSON.stringify({ tools, query, limit }));
    const inPtr = this.x.tb_alloc(input.length);
    // Re-read memory.buffer after every call that can grow memory.
    new Uint8Array(this.x.memory.buffer, inPtr, input.length).set(input);
    const packed = this.x.tb_search(inPtr, input.length);
    this.x.tb_free(inPtr, input.length);
    if (packed === 0n) throw new Error('wasm search failed (invalid input)');
    const outPtr = Number(packed >> 32n);
    const outLen = Number(packed & 0xffffffffn);
    const json = this.dec.decode(new Uint8Array(this.x.memory.buffer, outPtr, outLen));
    this.x.tb_free(outPtr, outLen);
    return JSON.parse(json) as string[];
  }
}
