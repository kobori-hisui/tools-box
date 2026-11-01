#!/usr/bin/env node
import { readdirSync } from 'node:fs';

const root = 'src/tools/categories';
let total = 0;
for (const cat of readdirSync(root, { withFileTypes: true }).filter((d) => d.isDirectory())) {
  const tools = readdirSync(`${root}/${cat.name}`, { withFileTypes: true }).filter((d) => d.isDirectory());
  total += tools.length;
  console.log(`${cat.name}: ${tools.map((t) => t.name).join(', ') || '-'}`);
}
console.log(`\nTotal: ${total}`);
