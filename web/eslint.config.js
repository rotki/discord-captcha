import rotki from '@rotki/eslint-config';

export default rotki({
  vue: true,
  typescript: {
    tsconfigPath: './tsconfig.json',
  },
  vueI18n: false,
  storybook: false,
  cypress: false,
  testLibrary: false,
});
