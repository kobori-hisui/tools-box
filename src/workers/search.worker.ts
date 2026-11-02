import { search, type Searchable } from '@/core/search/search';

/**
 * Off-main-thread search for very large catalogs.
 *   const w = new Worker(new URL('./search.worker.ts', import.meta.url), { type: 'module' });
 * Swap `search` for WasmSearch (core/wasm) once the catalog is big enough to matter.
 */
interface Req {
  id: number;
  tools: Searchable[];
  query: string;
  limit?: number;
}

globalThis.addEventListener('message', (e: Event) => {
  const { id, tools, query, limit } = (e as MessageEvent<Req>).data;
  const ids = search(tools, query, limit ?? 0).map((t) => t.id);
  (globalThis as unknown as { postMessage(m: unknown): void }).postMessage({ id, ids });
});
