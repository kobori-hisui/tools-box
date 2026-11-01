import { defineConfig } from 'vitest/config';
import { WxtVitest } from 'wxt/testing/vitest-plugin';

export default defineConfig({
  plugins: [WxtVitest()],
  test: {
    include: ['src/**/*.test.{ts,tsx}', 'tests/unit/**/*.test.ts'],
    setupFiles: ['tests/setup.ts'],
    mockReset: true,
    restoreMocks: true,
  },
});
