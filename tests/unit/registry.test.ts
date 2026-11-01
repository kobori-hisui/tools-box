import { describe, expect, it } from 'vitest';
import { getAllTools, searchTools } from '@/tools/_core/registry';

describe('registry', () => {
  it('descubre las herramientas de ejemplo', () => {
    expect(getAllTools().length).toBeGreaterThanOrEqual(3);
  });
  it('busca por keyword', () => {
    expect(searchTools('json').map((t) => t.id)).toContain('json-formatter');
  });
});
