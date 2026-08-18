import eslint from '@eslint/js';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import globals from 'globals';
import tseslint from 'typescript-eslint';

// typescript-eslint still needs the TypeScript 5/6 API. `tsc --noEmit` uses typescript-7 (TS 7).

export default tseslint.config(
  {
    ignores: [
      'public/**',
      'resources/**',
      'node_modules/**',
      'coverage/**',
      'api/**',
      'static/**',
    ],
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ['assets/js/**/*.{ts,tsx}'],
    plugins: {
      react,
      'react-hooks': reactHooks,
      'jsx-a11y': jsxA11y,
    },
    languageOptions: {
      globals: {
        ...globals.browser,
      },
      parserOptions: {
        ecmaFeatures: { jsx: true },
      },
    },
    settings: {
      react: { version: 'detect' },
    },
    rules: {
      ...react.configs.flat.recommended.rules,
      ...react.configs.flat['jsx-runtime'].rules,
      ...reactHooks.configs.recommended.rules,
      ...jsxA11y.flatConfigs.recommended.rules,
      // Islands fetch on mount; the compiler rule targets derived state, not data loading.
      'react-hooks/set-state-in-effect': 'off',
    },
  },
  {
    files: ['assets/js/**/*.test.{ts,tsx}', 'assets/js/test/**/*.ts'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.vitest,
      },
    },
  },
);
