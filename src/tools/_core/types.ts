import type { ContentScriptContext } from '#imports';
import type { ComponentType } from 'react';

export type ToolSurface =
  | 'popup' | 'sidepanel' | 'options' | 'content' | 'background' | 'devtools';

/** Metadata LIVIANA: se carga eager para buscar/listar herramientas. */
export interface ToolMeta {
  id: string;
  name: string;
  description: string;
  category: string;
  keywords: string[];
  icon: string;
  version: string;
  surfaces: ToolSurface[];
  /** Match patterns donde corre content.ts (omitir = nunca en paginas) */
  matches?: string[];
  /** Permisos opcionales que la herramienta pedira en runtime */
  permissions?: string[];
  browsers?: Array<'chrome' | 'firefox' | 'edge' | 'safari'>;
  /** Needs the Go host / Rust sidecar. Lists the protocol methods it calls. */
  native?: { methods: string[]; optional?: boolean };
  /** Extension APIs the tool needs (see platform/capabilities.ts CAPABILITIES keys). */
  requires?: string[];
  experimental?: boolean;
}

export interface ToolContentContext {
  ctx: ContentScriptContext;
  meta: ToolMeta;
  url: string;
}

export interface ToolBackgroundContext {
  meta: ToolMeta;
}

/** Implementacion PESADA: se carga lazy (import dinamico). */
export interface ToolImpl {
  Panel?: ComponentType;
  Options?: ComponentType;
  onBackground?(ctx: ToolBackgroundContext): void | Promise<void>;
  onContent?(ctx: ToolContentContext): void | Promise<void>;
}

export interface CategoryMeta {
  id: string;
  name: string;
  icon: string;
  order: number;
}
