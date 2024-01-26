import {
  type APIInteraction,
  ApplicationCommandsAPI,
  type Client,
  GatewayDispatchEvents,
  InteractionType,
  type WithIntrinsicProps,
} from '@discordjs/core';
import { commands } from '~/services/botcommands';
import { logger } from '~/utils/logger';
import { consume } from '~/utils/promise';
import type { REST } from '@discordjs/rest';

export class Commands {
  constructor(
    private readonly client: Client,
    private readonly rest: REST,
    private readonly config: { appId: string; guildId: string },
  ) {}

  async register() {
    const api = new ApplicationCommandsAPI(this.rest);

    for (const command of commands) {
      await api.createGuildCommand(
        this.config.appId,
        this.config.guildId,
        command.data,
      );
    }
    logger.info(`Registered ${commands.length} commands`);

    this.client.on(
      GatewayDispatchEvents.InteractionCreate,
      payload => consume(this.onInteractionCreate(payload)),
    );
  }

  private async onInteractionCreate({
    api,
    data: interaction,
  }: WithIntrinsicProps<APIInteraction>) {
    if (interaction.type !== InteractionType.ApplicationCommand) {
      logger.debug('Received interaction was not a command');
      return;
    }

    const currentCommand = commands.find(
      x => x.data.name === interaction.data.name,
    );
    if (currentCommand)
      await currentCommand.execute(interaction, api);
    else
      logger.debug(`Command ${interaction.data.name} is unknown`);
  }
}
