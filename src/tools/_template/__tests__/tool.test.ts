import { describe, expect, it } from 'vitest';
import meta from '../meta';

describe('__TOOL_ID__', () => {
  it('tiene metadata valida', () => {
    expect(meta.id).toBe('__TOOL_ID__');
    expect(meta.keywords.length).toBeGreaterThanOrEqual(0);
  });
});
