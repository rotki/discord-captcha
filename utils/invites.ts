import type {
  APIExtendedInvite,
  APIUser,
  GatewayInviteCreateDispatchData,
} from '@discordjs/core';
import type { CachedInvite, CachedUser } from '~/types/invites';

export const toCachedUser: {
  (user: APIUser): CachedUser;
  (user: APIUser | undefined): CachedUser | undefined;
} = (user: APIUser | undefined): any => {
  if (!user)
    return user;

  return {
    id: user.id,
    username: user.username,
  };
};

export function toCachedInvite(invite: APIExtendedInvite | GatewayInviteCreateDispatchData): CachedInvite {
  return {
    code: invite.code,
    data: {
      inviter: toCachedUser(invite.inviter),
      uses: invite.uses,
    },
  };
}
