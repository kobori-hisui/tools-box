import { describe, expect, it } from 'vitest';
import { CATALOG, support } from '@/entrypoints/onboarding/catalog';

describe('onboarding catalog', () => {
  it('has 100+ unique names', () => {
    expect(CATALOG.length).toBeGreaterThanOrEqual(100);
    expect(new Set(CATALOG.map((t) => t.name)).size).toBe(CATALOG.length);
  });
  it('every tool lists at least one api', () => {
    for (const t of CATALOG) expect(t.apis.length).toBeGreaterThan(0);
  });
  it('knows browser-specific apis', () => {
    expect(support('sidePanel')).toEqual({ chrome: true, firefox: false });
    expect(support('find')).toEqual({ chrome: false, firefox: true });
  });
});
