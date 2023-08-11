import { type API, type APIInteraction } from '@discordjs/core';

export const logsCommand = {
  data: {
    name: 'logsdir',
    description: 'Links to the documentation for the log files location',
  },
  async execute(interaction: APIInteraction, api: API) {
    await api.interactions.reply(interaction.id, interaction.token, {
      content:
        'You can find the default log file locations at https://rotki.readthedocs.io/en/stable/contribute.html#id2 ',
    });
  },
};
