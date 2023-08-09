import { REST } from '@discordjs/rest';
import { WebSocketManager } from '@discordjs/ws';
import { Client, GatewayIntentBits } from '@discordjs/core';
import { InviteMonitor } from '~/services/invite_monitor';

export default defineNuxtPlugin({
  name: 'discord',
  enforce: 'pre',
  async setup(nuxtApp) {
    const config = nuxtApp.$config.discord;
    const token = config.token;

    const rest = new REST({ version: '10' }).setToken(token);

    const gateway = new WebSocketManager({
      token,
      intents:
        GatewayIntentBits.Guilds |
        GatewayIntentBits.GuildInvites |
        GatewayIntentBits.GuildMembers,
      rest,
    });

    const client = new Client({ rest, gateway });
    const monitor = new InviteMonitor(client, {
      token,
      guildId: config.guildId.toString(),
      roleId: config.roleId.toString(),
    });
    monitor.setupListeners();
    await gateway.connect();

    return {
      provide: {
        discord: client,
      },
    };
  },
});
