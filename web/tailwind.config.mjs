import plugin from '@rotki/ui-library/theme';

export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts}',
  ],
  theme: {
    container: {
      center: true,
    },
  },
  plugins: [plugin],
};
