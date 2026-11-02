import type { ToolMeta } from '@/tools/_core/types';

/** Un solo mapa tipado de mensajes para toda la extension. */
export interface ProtocolMap {
  'tools:list': () => ToolMeta[];
  'tool:run': (data: { id: string; payload?: unknown }) => unknown;
  /** Forwarded to the Go native host (see core/native). */
  'native:call': (data: { method: string; params?: unknown }) => unknown;
  'native:available': () => boolean;
}
