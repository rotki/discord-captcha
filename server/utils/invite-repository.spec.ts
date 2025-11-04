import type { CachedInvite, CachedInviteData, CachedUser } from '../types/invites';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { InviteRepository } from './invite-repository';

const inviter: CachedUser = {
  id: '123456789',
  username: 'test',
};

function createInvite(code: string, data: Partial<CachedInviteData> = {}): CachedInvite {
  return {
    code,
    data: {
      expiresAt: new Date().toISOString(),
      inviter,
      maxUses: 1,
      uses: 0,
      ...data,
    },
  } satisfies CachedInvite;
}

describe('repository/InviteRepository', () => {
  let repo: InviteRepository;

  async function purge() {
    for await (const [code] of repo.iterator()) {
      await repo.delete(code);
    }
  }

  async function getInviteCodes(): Promise<string[]> {
    const codes: string[] = [];

    for await (const [code] of repo.iterator()) {
      codes.push(code);
    }
    return codes;
  }

  beforeEach(async () => {
    repo = new InviteRepository();
    vi.useFakeTimers();
    await purge();
  });

  afterEach(async () => {
    await purge();
  });

  it('should remove expired codes', async () => {
    const date = new Date();
    await repo.set(createInvite('0123', {
      expiresAt: date.toISOString(),
    }));
    vi.advanceTimersByTime(2000);

    await repo.set(createInvite('0124', {
      expiresAt: new Date().toISOString(),
    }));

    expect(await getInviteCodes()).toMatchObject(['0124']);
  });

  it('should not remove codes that never expire', async () => {
    await repo.set(createInvite('0123', {
      expiresAt: 'never',
    }));

    await repo.set(createInvite('0124'));

    expect(await getInviteCodes()).toMatchObject(['0123', '0124']);
  });

  it('should remove codes with undefined expiry', async () => {
    await repo.set(createInvite('0123', {
      expiresAt: undefined,
    }));

    await repo.set(createInvite('0124'));

    expect(await getInviteCodes()).toMatchObject(['0124']);
  });

  it('should remove codes that never expire but reached maxUses', async () => {
    await repo.set(createInvite('0123', {
      expiresAt: 'never',
      maxUses: 10,
      uses: 10,
    }));

    await repo.set(createInvite('0124'));

    expect(await getInviteCodes()).toMatchObject(['0124']);
  });
});
