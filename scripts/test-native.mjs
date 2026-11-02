#!/usr/bin/env node
// Runs the Rust and Go test suites. Rust goes first because the Go integration
// test talks to the real sidecar binary.
import { spawnSync } from 'node:child_process';
import { join, resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..');
const steps = [
  ['cargo', ['test', '--workspace'], 'native/rust'],
  ['cargo', ['clippy', '--all-targets', '--', '-D', 'warnings'], 'native/rust'],
  ['cargo', ['build', '--release', '-p', 'toolbox-sidecar'], 'native/rust'],
  ['go', ['vet', './...'], 'native/go'],
  ['go', ['test', '-race', '-count=1', './...'], 'native/go'],
];
for (const [cmd, args, cwd] of steps) {
  console.log(`\n$ (${cwd}) ${cmd} ${args.join(' ')}`);
  const r = spawnSync(cmd, args, { cwd: join(root, cwd), stdio: 'inherit' });
  if (r.status !== 0) process.exit(r.status ?? 1);
}
