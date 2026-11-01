import { defineConfig } from '@playwright/test';

// Ejecuta primero `pnpm build`. Los tests cargan .output/chrome-mv3
export default defineConfig({
  testDir: 'tests/e2e',
  fullyParallel: false,
  workers: 1,
  reporter: 'list',
});
