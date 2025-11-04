import plugin from '@rotki/ui-library/theme';

export default {
  content: [
    './app/components/**/*.{vue,js,ts}',
    './app/layouts/**/*.vue',
    './app/pages/**/*.vue',
    './app/app.vue',
    './app/error.vue',
  ],
  theme: {
    container: {
      center: true,
    },
  },
  plugins: [plugin],
};
