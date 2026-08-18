import { defineConfig } from 'vitest/config';
import path from 'node:path';

export default defineConfig({
  resolve: {
    alias: {
      '@params': path.resolve(__dirname, 'assets/js/test/params-stub.ts'),
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
