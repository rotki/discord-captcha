import { useRuntimeConfig } from '#imports';
import { createConsola } from 'consola';

export const logger = createConsola({
  fancy: true,
  level: useRuntimeConfig().public.logLevel,
});
