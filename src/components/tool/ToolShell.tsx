import { useMemo, useState } from 'react';
import { searchTools } from '@/tools/_core/registry';

export function ToolShell({ surface }: { surface: string }) {
  const [query, setQuery] = useState('');
  const results = useMemo(() => searchTools(query), [query]);

  return (
    <main className="tb-shell" data-surface={surface}>
      <input
        className="tb-search"
        placeholder="Search tools..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />
      <ul>
        {results.map((t) => (
          <li key={t.id}>{t.name}</li>
        ))}
      </ul>
    </main>
  );
}
