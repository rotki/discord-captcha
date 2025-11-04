import { logger } from './logger';

export function consume<T>(promise: Promise<T>): void {
  promise.then().catch(error => logger.error(error));
}
