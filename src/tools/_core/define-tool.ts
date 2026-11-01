import type { CategoryMeta, ToolImpl, ToolMeta } from './types';

export const defineToolMeta = (meta: ToolMeta): ToolMeta => meta;
export const defineTool = (impl: ToolImpl): ToolImpl => impl;
export const defineCategory = (meta: CategoryMeta): CategoryMeta => meta;
