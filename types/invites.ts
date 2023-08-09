export type CachedUser = { username: string; id: string };

export type CachedInvite = { code: string; data: CachedInviteData };

export type CachedInviteData = { uses: number; inviter?: CachedUser };
