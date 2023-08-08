import { FetchError } from 'ofetch';
import { $fetch } from 'ofetch/node';
import {
  CaptchaVerification,
  type DiscordInvite,
  DiscordInviteBody,
  DiscordInviteResponse,
} from '~/types/discord';
import { discordRequest } from '~/utils/discord';

export default defineEventHandler(async (event) => {
  const {
    recaptchaSecret,
    discord: { token, channelId },
  } = useRuntimeConfig();
  const requestBody = await readBody(event);
  const body = DiscordInviteBody.parse(requestBody);

  const response = CaptchaVerification.parse(
    await $fetch('https://www.google.com/recaptcha/api/siteverify', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams(
        Object.entries({
          secret: recaptchaSecret,
          response: body.captcha,
        }),
      ).toString(),
    }),
  );

  if (!response.success) {
    throw createError({
      statusCode: 400,
      statusMessage: response['error-codes']?.join(',') ?? '',
    });
  }

  try {
    const discordResponse = await discordRequest<DiscordInvite>(
      `/channels/${channelId}/invites`,
      {
        method: 'POST',
        // https://discord.com/developers/docs/resources/channel#create-channel-invite
        body: {
          max_age: 1800,
          max_uses: 1,
          unique: true,
        },
      },
      token,
    );

    return DiscordInviteResponse.parse(discordResponse);
  } catch (e: any) {
    if (e instanceof FetchError) {
      throw createError({
        statusCode: e.statusCode,
        statusMessage: `Invite creation failed with: ${e.statusMessage}`,
      });
    }

    throw createError({
      statusCode: 500,
      statusMessage: e.message,
    });
  }
});
