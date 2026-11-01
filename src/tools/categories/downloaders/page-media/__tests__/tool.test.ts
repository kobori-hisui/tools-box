import { describe, expect, it } from 'vitest';
import meta from '../meta';

describe('page-media', () => {
  it('tiene metadata valida', () => {
    expect(meta.id).toBe('page-media');
    expect(meta.keywords.length).toBeGreaterThanOrEqual(0);
  });
});
