import { REST } from '@discordjs/rest';
import { WebSocketManager } from '@discordjs/ws';
import { Client, GatewayIntentBits } from '@discordjs/core';
import { InviteMonitor } from '~/services/invite_monitor';
import { Commands } from '~/services/commands';
import { logger } from '~/utils/logger';

export default defineNitroPlugin(async () => {
  if (process.env.BUILD) {
    logger.info('Skipping plugin run due to BUILD mode');
    return;
  }
  const { discord: config } = useRuntimeConfig();
  const token = config.token;
  const guildId = config.guildId.toString();
  const roleId = config.roleId.toString();

  logger.info('Initializing bot gateway connection');

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
    guildId,
    roleId,
  });
  const commands = new Commands(client, rest, {
    appId: config.appId,
    guildId,
  });

  monitor.setupListeners();
  await commands.register();
  await gateway.connect();
});
