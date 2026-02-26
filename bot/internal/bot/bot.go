package bot

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/rotki/discord-captcha/internal/config"
	"github.com/rotki/discord-captcha/internal/store"
)

type Bot struct {
	session *discordgo.Session
	config  *config.Config
	store   store.InviteStore
	monitor *InviteMonitor
}

func New(cfg *config.Config, s store.InviteStore) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildInvites |
		discordgo.IntentsGuildMembers

	b := &Bot{
		session: session,
		config:  cfg,
		store:   s,
	}

	b.monitor = NewInviteMonitor(session, cfg, s)

	return b, nil
}

func (b *Bot) Start() error {
	b.monitor.Setup()
	registerCommands(b.session, b.config)

	if err := b.session.Open(); err != nil {
		return err
	}

	slog.Info("Discord bot connected")
	return nil
}

func (b *Bot) Stop() error {
	return b.session.Close()
}
