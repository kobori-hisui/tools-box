#!/usr/bin/env node
// macOS only. Wraps the Safari build in an Xcode project (requires Xcode).
// WXT's publishing guide documents this exact command:
//   pnpm wxt build -b safari && xcrun safari-web-extension-packager .output/safari-mv2
// Pass the .output/safari-mv2 (or safari-mv3) directory, never the source tree.
import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { join, resolve } from 'node:path';

if (process.platform !== 'darwin') {
  console.error('Safari packaging needs macOS + Xcode. Run this on a Mac.');
  process.exit(1);
}
const root = resolve(import.meta.dirname, '..');
const mv = process.argv.includes('--mv3') ? 'safari-mv3' : 'safari-mv2';

const build = spawnSync('pnpm', ['wxt', 'build', '-b', 'safari', mv === 'safari-mv3' ? '--mv3' : '--mv2'], {
  cwd: root,
  stdio: 'inherit',
});
if (build.status !== 0) process.exit(build.status ?? 1);

const out = join(root, '.output', mv);
if (!existsSync(out)) {
  console.error(`Missing ${out}`);
  process.exit(1);
}
// Extra converter flags (--app-name, --bundle-identifier, --project-location, ...):
// run `xcrun safari-web-extension-packager --help` and pass them after `--`.
const extra = process.argv.slice(2).filter((a) => a !== '--mv3');
const r = spawnSync('xcrun', ['safari-web-extension-packager', out, ...extra], { stdio: 'inherit' });
process.exit(r.status ?? 1);
