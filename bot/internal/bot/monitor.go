package bot

import (
	"log/slog"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rotki/discord-captcha/internal/config"
	"github.com/rotki/discord-captcha/internal/store"
)

type InviteMonitor struct {
	session *discordgo.Session
	config  *config.Config
	store   store.InviteStore
	mu      sync.RWMutex
	botUser store.CachedUser
}

func NewInviteMonitor(session *discordgo.Session, cfg *config.Config, s store.InviteStore) *InviteMonitor {
	return &InviteMonitor{
		session: session,
		config:  cfg,
		store:   s,
		botUser: store.CachedUser{ID: "", Username: "rotki"},
	}
}

func (m *InviteMonitor) getBotUserID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.botUser.ID
}

func (m *InviteMonitor) setBotUser(user store.CachedUser) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.botUser = user
}

func (m *InviteMonitor) Setup() {
	m.session.AddHandler(m.onReady)
	m.session.AddHandler(m.onGuildMemberAdd)
	m.session.AddHandler(m.onInviteCreate)
	m.session.AddHandler(m.onInviteDelete)
}

func (m *InviteMonitor) onReady(s *discordgo.Session, r *discordgo.Ready) {
	slog.Info("Gateway connection ready")
	time.Sleep(1 * time.Second)

	guildInvites, err := s.GuildInvites(m.config.DiscordGuildID)
	if err != nil {
		slog.Error("failed to fetch guild invites", "error", err)
		return
	}

	for _, invite := range guildInvites {
		if err := m.store.Set(toCachedInvite(invite)); err != nil {
			slog.Error("failed to cache invite", "code", invite.Code, "error", err)
		}
	}

	botUser := store.CachedUser{
		ID:       r.User.ID,
		Username: r.User.Username,
	}
	m.setBotUser(botUser)

	slog.Debug("invite monitor initialized",
		"bot", botUser.Username,
		"invites", len(guildInvites),
	)
}

func (m *InviteMonitor) onGuildMemberAdd(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	if event.User == nil {
		slog.Info("Missing user information, bailing.")
		return
	}

	member := store.CachedUser{
		ID:       event.User.ID,
		Username: event.User.Username,
	}

	slog.Debug("new user joined", "id", member.ID, "username", member.Username)

	apiInvites, err := s.GuildInvites(event.GuildID)
	if err != nil {
		slog.Error("failed to fetch guild invites on member add", "error", err)
		return
	}

	currentInvites := make(map[string]store.CachedInviteData)
	for _, invite := range apiInvites {
		cached := toCachedInvite(invite)
		currentInvites[cached.Code] = cached.Data
	}

	botUserID := m.getBotUserID()
	var usedBotInviteCode string

	for code, data := range m.store.Iterator() {
		invite, ok := currentInvites[code]
		if !ok {
			continue
		}

		if invite.Uses <= data.Uses {
			continue
		}

		// Always update the store with fresh use count
		if err := m.store.Set(store.CachedInvite{Code: code, Data: invite}); err != nil {
			slog.Error("failed to update invite", "code", code, "error", err)
		}

		if invite.Inviter != nil && invite.Inviter.ID != botUserID {
			slog.Debug("invite used by non-bot inviter, bailing", "code", code)
			return
		}

		usedBotInviteCode = code
	}

	if usedBotInviteCode == "" {
		slog.Warn("could not determine which invite was used", "user", member.Username)
		return
	}

	roleID := m.config.DiscordRoleID
	slog.Debug("adding role to user", "username", member.Username, "role", roleID, "invite", usedBotInviteCode)

	if err := s.GuildMemberRoleAdd(event.GuildID, member.ID, roleID); err != nil {
		slog.Error("failed to add role", "user", member.Username, "role", roleID, "error", err)
	}
}

func (m *InviteMonitor) onInviteCreate(s *discordgo.Session, event *discordgo.InviteCreate) {
	inviterName := ""
	if event.Inviter != nil {
		inviterName = event.Inviter.Username
	}
	slog.Info("invite created", "code", event.Code, "inviter", inviterName)

	if err := m.store.Set(toCachedInviteFromCreate(event)); err != nil {
		slog.Error("failed to cache new invite", "code", event.Code, "error", err)
	}
}

func (m *InviteMonitor) onInviteDelete(s *discordgo.Session, event *discordgo.InviteDelete) {
	slog.Info("invite deleted", "code", event.Code)

	if err := m.store.Delete(event.Code); err != nil {
		slog.Error("failed to delete invite from cache", "code", event.Code, "error", err)
	}
}

func toCachedInvite(invite *discordgo.Invite) store.CachedInvite {
	expiresAt := "never"
	if invite.ExpiresAt != nil {
		expiresAt = invite.ExpiresAt.Format(time.RFC3339)
	}

	var inviter *store.CachedUser
	if invite.Inviter != nil {
		inviter = &store.CachedUser{
			ID:       invite.Inviter.ID,
			Username: invite.Inviter.Username,
		}
	}

	return store.CachedInvite{
		Code: invite.Code,
		Data: store.CachedInviteData{
			Uses:      invite.Uses,
			MaxUses:   invite.MaxUses,
			Inviter:   inviter,
			ExpiresAt: expiresAt,
		},
	}
}

func toCachedInviteFromCreate(event *discordgo.InviteCreate) store.CachedInvite {
	return toCachedInvite(event.Invite)
}
