#!/usr/bin/env node
// Builds the Go binaries (all three) and the Rust sidecar into native/dist/<os>-<arch>/.
//   pnpm native:build                      -> current platform
//   pnpm native:build -- --all             -> cross-compile Go for win/mac/linux
//   pnpm native:build:wasm                 -> also build public/wasm/toolbox.wasm
import { spawnSync } from 'node:child_process';
import { copyFileSync, existsSync, mkdirSync } from 'node:fs';
import { join, resolve } from 'node:path';

const root = resolve(import.meta.dirname, '..');
const args = new Set(process.argv.slice(2));
const run = (cmd, argv, opts = {}) => {
  const r = spawnSync(cmd, argv, { stdio: 'inherit', ...opts });
  if (r.status !== 0) {
    console.error(`\n${cmd} ${argv.join(' ')} failed`);
    process.exit(r.status ?? 1);
  }
};

const goos = { linux: 'linux', darwin: 'darwin', win32: 'windows' }[process.platform];
const goarch = { x64: 'amd64', arm64: 'arm64' }[process.arch];
const targets = args.has('--all')
  ? [['linux', 'amd64'], ['linux', 'arm64'], ['darwin', 'amd64'], ['darwin', 'arm64'], ['windows', 'amd64'], ['windows', 'arm64']]
  : [[goos, goarch]];

for (const [os, arch] of targets) {
  const out = join(root, 'native/dist', `${os}-${arch}`);
  mkdirSync(out, { recursive: true });
  const ext = os === 'windows' ? '.exe' : '';
  for (const cmd of ['toolbox-host', 'toolboxd', 'toolbox-installer']) {
    console.log(`go build ${cmd} (${os}/${arch})`);
    run('go', ['build', '-trimpath', '-ldflags', '-s -w', '-o', join(out, cmd + ext), `./cmd/${cmd}`], {
      cwd: join(root, 'native/go'),
      env: { ...process.env, GOOS: os, GOARCH: arch, CGO_ENABLED: '0' },
    });
  }
}

// The Rust sidecar is built for the host platform only (cross-compiling needs per-target linkers).
console.log('cargo build --release -p toolbox-sidecar');
run('cargo', ['build', '--release', '-p', 'toolbox-sidecar'], { cwd: join(root, 'native/rust') });
const sidecarExt = process.platform === 'win32' ? '.exe' : '';
copyFileSync(
  join(root, 'native/rust/target/release', `toolbox-sidecar${sidecarExt}`),
  join(root, 'native/dist', `${goos}-${goarch}`, `toolbox-sidecar${sidecarExt}`),
);

if (args.has('--wasm')) {
  console.log('cargo build --profile wasm -p toolbox-wasm --target wasm32-unknown-unknown');
  run('rustup', ['target', 'add', 'wasm32-unknown-unknown']);
  run('cargo', ['build', '--profile', 'wasm', '-p', 'toolbox-wasm', '--target', 'wasm32-unknown-unknown'], {
    cwd: join(root, 'native/rust'),
  });
  const wasm = join(root, 'native/rust/target/wasm32-unknown-unknown/wasm/toolbox_wasm.wasm');
  if (!existsSync(wasm)) {
    console.error('wasm output not found: ' + wasm);
    process.exit(1);
  }
  mkdirSync(join(root, 'public/wasm'), { recursive: true });
  copyFileSync(wasm, join(root, 'public/wasm/toolbox.wasm'));
  console.log('-> public/wasm/toolbox.wasm');
}
console.log('\nDone: native/dist');
