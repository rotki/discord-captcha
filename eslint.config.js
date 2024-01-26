import rotki from '@rotki/eslint-config';
import pluginNuxt from 'eslint-plugin-nuxt';

export default rotki({
  vue: true,
  typescript: {
    tsconfigPath: './.nuxt/tsconfig.json',
  },
  stylistic: true,
  formatters: true,
}, {
  files: ['**/*.ts'],
  rules: {
    'perfectionist/sort-objects': 'error',
  },
}, {
  plugins: {
    nuxt: pluginNuxt,
  },
  rules: {
    ...pluginNuxt.configs.base.rules,
    ...pluginNuxt.configs.recommended.rules,
  },
}, {
  files: [
    'modules/**/*.ts',
    'server/**/*.ts',
    'types/recaptcha.d.ts',
  ],
  rules: {
    'import/no-default-export': 'off',
  },
});
