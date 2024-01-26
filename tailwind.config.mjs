import plugin from '@rotki/ui-library/theme';

export default {
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
  ],
  theme: {
    container: {
      center: true,
    },
  },
  plugins: [plugin],
};
