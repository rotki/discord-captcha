import { defineVitestConfig } from '@nuxt/test-utils/config';

export default defineVitestConfig({
  test: {
    coverage: {
      exclude: ['node_modules', '**/*.d.ts', '**/*.spec.ts', '**/*.config.ts'],
      include: ['**/*.ts', '**/*.vue'],
      provider: 'v8',
      reporter: ['json'],
    },
    environment: 'nuxt',
  },
});
