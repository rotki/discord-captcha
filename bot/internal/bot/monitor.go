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

	// Phase 1: Snapshot — read cached and live state into plain maps
	cachedInvites := snapshotStore(m.store)

	apiInvites, err := s.GuildInvites(event.GuildID)
	if err != nil {
		slog.Error("failed to fetch guild invites on member add", "error", err)
		return
	}

	currentInvites := make(map[string]store.CachedInviteData, len(apiInvites))
	for _, invite := range apiInvites {
		cached := toCachedInvite(invite)
		currentInvites[cached.Code] = cached.Data
	}

	// Phase 2: Match — pure function, no side effects
	result := matchUsedInvite(cachedInvites, currentInvites, m.getBotUserID())

	// Phase 3: Mutate — apply store updates outside any iteration
	for code, data := range result.StoreUpdates {
		if err := m.store.Set(store.CachedInvite{Code: code, Data: data}); err != nil {
			slog.Error("failed to update invite", "code", code, "error", err)
		}
	}

	if result.Code == "" {
		slog.Warn("could not determine which invite was used", "user", member.Username)
		return
	}

	if err := m.store.Delete(result.Code); err != nil {
		slog.Error("failed to delete consumed invite from cache", "code", result.Code, "error", err)
	}

	roleID := m.config.DiscordRoleID
	slog.Debug("adding role to user", "username", member.Username, "role", roleID, "invite", result.Code)

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

func (m *InviteMonitor) onInviteDelete(_ *discordgo.Session, event *discordgo.InviteDelete) {
	// Intentionally no cache deletion here. Discord fires InviteDelete before
	// GuildMemberAdd for single-use invites, and onGuildMemberAdd needs the
	// cached entry to detect consumed invites. Cleanup is handled by
	// onGuildMemberAdd (after match) and the periodic Cleanup goroutine.
	slog.Info("invite deleted", "code", event.Code)
}

func snapshotStore(s store.InviteStore) map[string]store.CachedInviteData {
	snap := make(map[string]store.CachedInviteData)
	for code, data := range s.Iterator() {
		snap[code] = data
	}
	return snap
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
