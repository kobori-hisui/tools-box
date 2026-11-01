import { describe, expect, it } from 'vitest';
import meta from '../meta';

describe('json-formatter', () => {
  it('tiene metadata valida', () => {
    expect(meta.id).toBe('json-formatter');
    expect(meta.keywords.length).toBeGreaterThanOrEqual(0);
  });
});
