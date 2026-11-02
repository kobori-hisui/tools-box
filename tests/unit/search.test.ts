import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { describe, expect, it } from 'vitest';
import { search, words, type Searchable } from '@/core/search/search';

interface Fixture {
  tools: Searchable[];
  cases: Array<{ query: string; limit?: number; expected: string[] }>;
}

const fixture: Fixture = JSON.parse(
  readFileSync(resolve(__dirname, '../../protocol/fixtures/search.cases.json'), 'utf8'),
);

describe('search (parity with the Rust implementation)', () => {
  for (const c of fixture.cases) {
    it(`"${c.query}"${c.limit ? ` limit ${c.limit}` : ''}`, () => {
      expect(search(fixture.tools, c.query, c.limit ?? 0).map((t) => t.id)).toEqual(c.expected);
    });
  }

  it('splits on punctuation and lowercases', () => {
    expect(words('images/videos, JSON.')).toEqual(['images', 'videos', 'json']);
  });

  it('does not mutate the input order', () => {
    const ids = fixture.tools.map((t) => t.id);
    search(fixture.tools, 'encode');
    expect(fixture.tools.map((t) => t.id)).toEqual(ids);
  });
});
