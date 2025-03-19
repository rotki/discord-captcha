import type { CachedInvite, CachedInviteData } from '~/types/invites';
import { createStorage, type Driver, type Storage } from 'unstorage';
import fsDriver from 'unstorage/drivers/fs';
import redisDriver from 'unstorage/drivers/redis';
import { logger } from '~/utils/logger';

export class InviteRepository {
  private readonly storage: Storage<CachedInviteData>;

  constructor() {
    const { redis } = useRuntimeConfig();

    let driver: Driver;
    if (redis.host && redis.password !== undefined) {
      logger.info('Using unstorage redis driver');
      driver = redisDriver({
        base: 'discord_invites',
        host: redis.host,
        password: redis.password,
      });
    }
    else {
      logger.info('Using unstorage fs driver');
      driver = fsDriver({ base: './data/invites' });
    }
    this.storage = createStorage({
      driver,
    });
  }

  async set(cachedInvite: CachedInvite): Promise<void> {
    await this.storage.setItem(cachedInvite.code, cachedInvite.data);
    await this.cleanup();
  }

  private async cleanup(): Promise<void> {
    const now = Date.now();
    for await (const [code, data] of this.iterator()) {
      const expiresAt = data.expiresAt;
      if (expiresAt === 'never') {
        if (data.maxUses > 0 && data.maxUses === data.uses) {
          logger.info(`invite ${code} reached max uses, purging`);
          await this.delete(code);
        }
      }
      else if (expiresAt === undefined) {
        logger.info(`invite ${code} didn't have expiration data, purging`);
        await this.delete(code);
      }
      else {
        const expirationDate = new Date(expiresAt);
        if (expirationDate.getTime() < now) {
          logger.info(`invite ${code} expired at ${expiresAt}, purging`);
          await this.delete(code);
        }
      }
    }
  }

  async delete(code: string): Promise<void> {
    await this.storage.removeItem(code);
  }

  async* iterator(): AsyncGenerator<[string, CachedInviteData]> {
    const keys = await this.storage.getKeys();
    for (const key of keys) {
      const value = await this.storage.getItem(key);
      if (value)
        yield [key, value];
    }
  }
}
