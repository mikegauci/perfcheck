/** @type {import('stylelint').Config} */
export default {
  extends: ['stylelint-config-standard-scss'],
  ignoreFiles: ['public/**', 'node_modules/**', 'resources/**'],
  rules: {
    'selector-class-pattern': [
      '^perfcheck(-[a-z0-9]+)+(__[a-z0-9]+(-[a-z0-9]+)*)?(--[a-z0-9]+(-[a-z0-9]+)*)?$',
      {
        resolveNestedSelectors: true,
        message: 'Expected class to follow BEM: perfcheck-block__elem--mod',
      },
    ],
    'no-descending-specificity': null,
    'scss/load-partial-extension': null,
    'scss/at-use-no-unnamespaced': null,
    'scss/dollar-variable-empty-line-before': null,
    'at-rule-empty-line-before': null,
    'media-feature-range-notation': null,
    'color-hex-length': 'long',
    'value-keyword-case': null,
    'property-no-vendor-prefix': null,
    'property-no-deprecated': [true, { ignoreProperties: ['clip'] }],
    'declaration-property-value-keyword-no-deprecated': [
      true,
      { ignoreKeywords: ['break-word'] },
    ],
  },
};

