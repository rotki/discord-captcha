export interface CachedUser { username: string; id: string }

export interface CachedInvite { code: string; data: CachedInviteData }

export interface CachedInviteData { uses: number; inviter?: CachedUser }
