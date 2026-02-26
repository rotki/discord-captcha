package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/rotki/discord-captcha/internal/config"
)

var slashCommands = []*discordgo.ApplicationCommand{
	{
		Name:        "logsdir",
		Description: "Links to the documentation for the log files location",
	},
	{
		Name:        "datadir",
		Description: "Links to the documentation for the data directory location",
	},
}

var commandResponses = map[string]string{
	"logsdir": "You can find the default log file locations at https://docs.rotki.com/contribution-guides/#run-rotki-in-debug-mode",
	"datadir": "You can find the default data directory locations at https://docs.rotki.com/usage-guides/data-directory.html#rotki-data-directory",
}

func registerCommands(s *discordgo.Session, cfg *config.Config) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		response, ok := commandResponses[i.ApplicationCommandData().Name]
		if !ok {
			slog.Debug("unknown command", "name", i.ApplicationCommandData().Name)
			return
		}

		if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		}); err != nil {
			slog.Error("failed to respond to command", "name", i.ApplicationCommandData().Name, "error", err)
		}
	})

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		for _, cmd := range slashCommands {
			if _, err := s.ApplicationCommandCreate(cfg.DiscordAppID, cfg.DiscordGuildID, cmd); err != nil {
				slog.Error("failed to register command", "name", cmd.Name, "error", err)
			}
		}
		slog.Info("registered slash commands", "count", len(slashCommands))
	})
}
