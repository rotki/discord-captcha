import {
  type APIExtendedInvite,
  type APIUser,
  type GatewayInviteCreateDispatchData,
} from '@discordjs/core';
import { type CachedInvite, type CachedUser } from '~/types/invites';

export const toCachedUser: {
  (user: APIUser): CachedUser;
  (user: APIUser | undefined): CachedUser | undefined;
} = (user: APIUser | undefined): any => {
  if (!user) {
    return user;
  }

  return {
    id: user.id,
    username: user.username,
  };
};

export const toCachedInvite = (
  invite: APIExtendedInvite | GatewayInviteCreateDispatchData,
): CachedInvite => ({
  code: invite.code,
  data: {
    uses: invite.uses,
    inviter: toCachedUser(invite.inviter),
  },
});
