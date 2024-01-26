import { $fetch, type FetchOptions, type FetchRequest } from 'ofetch/node';

export async function discordRequest<
  Resp,
  Req extends FetchRequest = FetchRequest,
>(endpoint: Req, options: FetchOptions<'json'>, token: string): Promise<Resp> {
  const url = `https://discord.com/api/v10${endpoint.toString()}`;
  if (options.body)
    options.body = JSON.stringify(options.body);

  return await $fetch<Resp>(url, {
    headers: {
      'Authorization': `Bot ${token}`,
      'Content-Type': 'application/json; charset=UTF-8',
      'User-Agent': 'RotkiBot',
    },
    ...options,
  });
}
