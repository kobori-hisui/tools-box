import { describe, expect, it } from 'vitest';
import meta from '../meta';

describe('word-counter', () => {
  it('tiene metadata valida', () => {
    expect(meta.id).toBe('word-counter');
    expect(meta.keywords.length).toBeGreaterThanOrEqual(0);
  });
});
