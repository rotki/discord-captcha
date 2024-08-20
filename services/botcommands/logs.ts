import type { API, APIInteraction } from '@discordjs/core';

export const logsCommand = {
  data: {
    description: 'Links to the documentation for the log files location',
    name: 'logsdir',
  },
  async execute(interaction: APIInteraction, api: API) {
    await api.interactions.reply(interaction.id, interaction.token, {
      content:
        'You can find the default log file locations at https://docs.rotki.com/contribution-guides/#run-rotki-in-debug-mode',
    });
  },
};
