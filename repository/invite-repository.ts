import { type Driver, type Storage, createStorage } from 'unstorage';
import fsDriver from 'unstorage/drivers/fs';
import redisDriver from 'unstorage/drivers/redis';
import { logger } from '~/utils/logger';
import type { CachedInvite, CachedInviteData } from '~/types/invites';

export class InviteRepository {
  private readonly storage: Storage<CachedInviteData>;

  constructor() {
    const { redis } = useRuntimeConfig();

    let driver: Driver;
    if (redis.host && redis.password) {
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

  async set(cachedInvite: CachedInvite) {
    await this.storage.setItem(cachedInvite.code, cachedInvite.data);
  }

  async delete(code: string): Promise<void> {
    await this.storage.removeItem(code);
  }

  async *iterator(): AsyncGenerator<[string, CachedInviteData]> {
    const keys = await this.storage.getKeys();
    for (const key of keys) {
      const value = await this.storage.getItem(key);
      if (value)
        yield [key, value];
    }
  }
}
