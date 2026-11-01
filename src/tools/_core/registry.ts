import { MatchPattern } from 'wxt/utils/match-patterns';
import type { CategoryMeta, ToolImpl, ToolMeta } from './types';

// meta.ts: eager (liviano). index.ts: lazy (code-splitting por herramienta).
const metaModules = import.meta.glob<{ default: ToolMeta }>(
  '../categories/*/*/meta.ts',
  { eager: true },
);
const implLoaders = import.meta.glob<{ default: ToolImpl }>('../categories/*/*/index.ts');
const categoryModules = import.meta.glob<{ default: CategoryMeta }>(
  '../categories/*/_category.ts',
  { eager: true },
);

const tools: ToolMeta[] = Object.values(metaModules).map((m) => m.default);
const byId = new Map(tools.map((t) => [t.id, t]));
const implKeyById = new Map(
  Object.keys(metaModules).map((key) => [
    metaModules[key]!.default.id,
    key.replace(/meta\.ts$/, 'index.ts'),
  ]),
);

export const getAllTools = (): ToolMeta[] => tools;
export const getTool = (id: string) => byId.get(id);
export const getCategories = (): CategoryMeta[] =>
  Object.values(categoryModules).map((m) => m.default).sort((a, b) => a.order - b.order);

export function searchTools(query: string): ToolMeta[] {
  const q = query.trim().toLowerCase();
  if (!q) return tools;
  return tools.filter((t) =>
    [t.name, t.description, t.category, ...t.keywords].join(' ').toLowerCase().includes(q),
  );
}

export function getToolsForUrl(url: string): ToolMeta[] {
  return tools.filter((t) => t.matches?.some((m) => new MatchPattern(m).includes(url)));
}

export async function loadToolImpl(id: string): Promise<ToolImpl | undefined> {
  const key = implKeyById.get(id);
  const loader = key ? implLoaders[key] : undefined;
  return (await loader?.())?.default;
}
