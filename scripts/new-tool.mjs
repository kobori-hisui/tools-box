#!/usr/bin/env node
// uso: pnpm tool:new <categoria> <tool-id> ["Nombre visible"]
import { cp, readdir, readFile, writeFile, stat } from 'node:fs/promises';
import { join } from 'node:path';

const [category, id, ...nameParts] = process.argv.slice(2);
if (!category || !id) {
  console.error('uso: pnpm tool:new <categoria> <tool-id> ["Nombre"]');
  process.exit(1);
}
const name = nameParts.join(' ') || id.replace(/-/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
const src = 'src/tools/_template';
const dest = `src/tools/categories/${category}/${id}`;

async function walk(dir) {
  for (const entry of await readdir(dir)) {
    const p = join(dir, entry);
    if ((await stat(p)).isDirectory()) await walk(p);
    else {
      const text = await readFile(p, 'utf8');
      await writeFile(
        p,
        text
          .replaceAll('__TOOL_ID__', id)
          .replaceAll('__TOOL_NAME__', name)
          .replaceAll('__TOOL_DESCRIPTION__', `${name}.`)
          .replaceAll('__CATEGORY__', category),
      );
    }
  }
}

await cp(src, dest, { recursive: true, errorOnExist: true });
await walk(dest);
console.log(`Herramienta creada en ${dest}`);
