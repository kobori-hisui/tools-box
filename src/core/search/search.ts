/**
 * Tool search. Mirrors native/rust/crates/core/src/search.rs; both are tested
 * against /protocol/fixtures/search.cases.json, so they cannot drift apart.
 *
 * Every query token must match (AND). A token scores the best of:
 *   name word exact 100 | keyword exact 90 | name word prefix 80 |
 *   keyword prefix 70   | category exact 40 | description word prefix 30
 * Results are sorted by total score (descending); ties keep input order.
 */
export interface Searchable {
  id: string;
  name: string;
  description?: string;
  category?: string;
  keywords?: readonly string[];
}

export function words(s: string): string[] {
  return s
    .toLowerCase()
    .split(/[^\p{L}\p{N}]+/u)
    .filter(Boolean);
}

function bestIn(token: string, ws: string[], exact: number, prefix: number): number {
  let best = 0;
  for (const w of ws) {
    if (w === token) return exact;
    if (w.startsWith(token)) best = Math.max(best, prefix);
  }
  return best;
}

interface Prepared {
  name: string[];
  keywords: string[];
  category: string[];
  description: string[];
}

function prepare(t: Searchable): Prepared {
  return {
    name: words(t.name),
    keywords: (t.keywords ?? []).flatMap(words),
    category: words(t.category ?? ''),
    description: words(t.description ?? ''),
  };
}

function scoreToken(p: Prepared, token: string): number {
  let best = bestIn(token, p.name, 100, 80);
  best = Math.max(best, bestIn(token, p.keywords, 90, 70));
  if (p.category.includes(token)) best = Math.max(best, 40);
  if (p.description.some((w) => w.startsWith(token))) best = Math.max(best, 30);
  return best;
}

/** `limit <= 0` means no limit. */
export function search<T extends Searchable>(tools: readonly T[], query: string, limit = 0): T[] {
  const tokens = words(query);
  const scored: Array<{ score: number; index: number }> = [];
  tools.forEach((t, index) => {
    if (tokens.length === 0) {
      scored.push({ score: 0, index });
      return;
    }
    const p = prepare(t);
    let total = 0;
    for (const tok of tokens) {
      const s = scoreToken(p, tok);
      if (s === 0) return;
      total += s;
    }
    scored.push({ score: total, index });
  });
  scored.sort((a, b) => b.score - a.score || a.index - b.index);
  const out = scored.map((s) => tools[s.index]!);
  return limit > 0 ? out.slice(0, limit) : out;
}
