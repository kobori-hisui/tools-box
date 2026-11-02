import { describe, expect, it } from 'vitest';
import { CAPABILITIES, has, snapshot, supports } from '@/platform/capabilities';

describe('capabilities', () => {
  it('has() walks dotted paths and never throws', () => {
    expect(has('definitely.not.there')).toBe(false);
    expect(has('')).toBe(false);
  });

  it('snapshot covers every declared capability with booleans', () => {
    const snap = snapshot();
    expect(Object.keys(snap).sort()).toEqual(Object.keys(CAPABILITIES).sort());
    expect(Object.values(snap).every((v) => typeof v === 'boolean')).toBe(true);
  });

  it('supports() agrees with has()', () => {
    expect(supports('debugger')).toBe(has(CAPABILITIES.debugger));
  });
});
