import { useMemo, useState } from 'react';
import './style.css';
import { CATALOG, support } from './catalog';

export default function App() {
  const [q, setQ] = useState('');
  const groups = useMemo(() => {
    const s = q.trim().toLowerCase();
    const map = new Map<string, typeof CATALOG>();
    for (const t of CATALOG) {
      if (
        s &&
        !`${t.name} ${t.category} ${t.apis.join(' ')}`.toLowerCase().includes(s)
      )
        continue;
      map.set(t.category, [...(map.get(t.category) ?? []), t]);
    }
    return [...map.entries()];
  }, [q]);

  return (
    <main className="onb">
      <h1>Toolbox · {CATALOG.length} ideas de herramientas</h1>
      <p className="onb-note">
        Solo nombres. Chips: APIs que usaria (C = Chrome, F = Firefox, Safari
        sin verificar).
      </p>
      <input
        value={q}
        onChange={(e) => setQ(e.target.value)}
        placeholder="Buscar nombre, categoria o API"
      />
      {groups.map(([cat, tools]) => (
        <section key={cat}>
          <h2>
            {cat} <small>({tools.length})</small>
          </h2>
          <ul>
            {tools.map((t) => (
              <li key={t.name}>
                <strong>{t.name}</strong>
                <span className="apis">
                  {t.apis.map((a) => {
                    const s = support(a);
                    return (
                      <code
                        key={a}
                        title={`Chrome: ${s.chrome ? 'si' : 'no'} · Firefox: ${s.firefox ? 'si' : 'no'} · Safari: ?`}
                      >
                        {a} {s.chrome ? 'C' : ''}
                        {s.firefox ? 'F' : ''}
                      </code>
                    );
                  })}
                </span>
              </li>
            ))}
          </ul>
        </section>
      ))}
    </main>
  );
}
