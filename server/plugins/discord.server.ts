import process from 'node:process';
import { REST } from '@discordjs/rest';
import { WebSocketManager } from '@discordjs/ws';
import { Client, GatewayIntentBits } from '@discordjs/core';
import { InviteMonitor } from '~/services/invite-monitor';
import { Commands } from '~/services/commands';
import { logger } from '~/utils/logger';
import { consume } from '~/utils/promise';

async function initPlugin() {
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
    intents:
      GatewayIntentBits.Guilds
      | GatewayIntentBits.GuildInvites
      | GatewayIntentBits.GuildMembers,
    rest,
    token,
  });

  const client = new Client({ gateway, rest });

  const monitor = new InviteMonitor(client, {
    guildId,
    roleId,
    token,
  });
  const commands = new Commands(client, rest, {
    appId: config.appId,
    guildId,
  });

  monitor.setupListeners();
  await commands.register();
  await gateway.connect();
}

export default defineNitroPlugin(() => consume(initPlugin()));
