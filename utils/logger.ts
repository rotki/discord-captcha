import { createConsola } from 'consola';

export const logger = createConsola({
  level: useRuntimeConfig().public.logLevel,
  fancy: true,
});
