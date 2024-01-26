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
    discord: { channelId, token },
    recaptchaSecret,
  } = useRuntimeConfig();
  const requestBody = await readBody(event);
  const body = DiscordInviteBody.parse(requestBody);

  const response = CaptchaVerification.parse(
    await $fetch('https://www.google.com/recaptcha/api/siteverify', {
      body: new URLSearchParams(
        Object.entries({
          response: body.captcha,
          secret: recaptchaSecret,
        }),
      ).toString(),
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      method: 'POST',
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
        // https://discord.com/developers/docs/resources/channel#create-channel-invite
        body: {
          max_age: 1800,
          max_uses: 1,
          unique: true,
        },
        method: 'POST',
      },
      token,
    );

    return DiscordInviteResponse.parse(discordResponse);
  }
  catch (error: any) {
    if (error instanceof FetchError) {
      throw createError({
        statusCode: error.statusCode,
        statusMessage: `Invite creation failed with: ${error.statusMessage}`,
      });
    }

    throw createError({
      statusCode: 500,
      statusMessage: error.message,
    });
  }
});
