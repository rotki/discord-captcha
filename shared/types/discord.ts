import { z } from 'zod';

export const DiscordInviteBody = z.object({
  captcha: z.string(),
});

export const CaptchaVerification = z.object({
  'challenge_ts': z.string().optional(),
  'error-codes': z.array(z.string()).optional(),
  'hostname': z.string().optional(),
  'success': z.boolean(),
});

export type CaptchaVerification = z.infer<typeof CaptchaVerification>;

export const DiscordInviteResponse = z.object({
  code: z.string(),
  expires_at: z.coerce.date(),
});

export type DiscordInvite = z.infer<typeof DiscordInviteResponse>;
