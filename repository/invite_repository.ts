import { type Storage, createStorage } from 'unstorage';
import fsDriver from 'unstorage/drivers/fs';
import { type CachedInvite, type CachedInviteData } from '~/types/invites';

export class InviteRepository {
  private readonly storage: Storage<CachedInviteData>;

  constructor() {
    this.storage = createStorage({
      driver: fsDriver({ base: './data/invites' }),
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
      if (value) {
        yield [key, value];
      }
    }
  }
}
