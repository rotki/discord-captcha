import { logger } from '~/utils/logger';

export function consume<T>(promise: Promise<T>): void {
  promise.then().catch(error => logger.error(error));
}
