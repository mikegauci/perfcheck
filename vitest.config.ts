import { defineConfig } from 'vitest/config';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  resolve: {
    alias: {
      '@params': path.resolve(root, 'assets/js/test/params-stub.ts'),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./assets/js/test/setup.ts'],
    include: ['assets/js/**/*.{test,spec}.{ts,tsx}'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html'],
      include: ['assets/js/**/*.{ts,tsx}'],
      exclude: ['assets/js/test/**', 'assets/js/**/*.test.*', 'assets/js/main.tsx'],
    },
  },
});
