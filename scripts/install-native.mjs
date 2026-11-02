#!/usr/bin/env node
// Registers the native host with the browsers on this machine.
//   pnpm native:install -- --chrome-id <ID> --firefox-id toolbox@example.com
// Chrome's extension id is stable only if you pin `key` in the manifest (or after publishing).
import { spawnSync } from 'node:child_process';
import { existsSync } from 'node:fs';
import { join, resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..');
const goos = { linux: 'linux', darwin: 'darwin', win32: 'windows' }[process.platform];
const goarch = { x64: 'amd64', arm64: 'arm64' }[process.arch];
const ext = goos === 'windows' ? '.exe' : '';
const dist = join(root, 'native/dist', `${goos}-${goarch}`);
const installer = join(dist, `toolbox-installer${ext}`);
const host = join(dist, `toolbox-host${ext}`);

if (!existsSync(installer) || !existsSync(host)) {
  console.error('Binaries not found. Run `pnpm native:build` first.');
  process.exit(1);
}
const mode = process.argv.includes('--uninstall') ? 'uninstall' : process.argv.includes('--print') ? 'print' : 'install';
const passthrough = process.argv.slice(2).filter((a) => !['--uninstall', '--print', '--'].includes(a));
const r = spawnSync(installer, [mode, '--host', host, ...passthrough], { stdio: 'inherit' });
process.exit(r.status ?? 1);
