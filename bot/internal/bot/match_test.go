package bot

import (
	"testing"

	"github.com/rotki/discord-captcha/internal/store"
)

func botUser() *store.CachedUser {
	return &store.CachedUser{ID: "bot-1", Username: "rotki"}
}

func otherUser() *store.CachedUser {
	return &store.CachedUser{ID: "other-1", Username: "someone"}
}

func inv(uses, maxUses int, inviter *store.CachedUser) store.CachedInviteData {
	return store.CachedInviteData{
		Uses:      uses,
		MaxUses:   maxUses,
		Inviter:   inviter,
		ExpiresAt: "never",
	}
}

func TestMatchUsedInvite(t *testing.T) {
	const botID = "bot-1"

	tests := []struct {
		name             string
		cached           map[string]store.CachedInviteData
		api              map[string]store.CachedInviteData
		wantCode         string
		wantUpdateCodes  []string // codes expected in StoreUpdates
		wantNoUpdateCode []string // codes expected NOT in StoreUpdates
	}{
		{
			name:            "use count increased on bot invite",
			cached:          map[string]store.CachedInviteData{"abc": inv(0, 5, botUser())},
			api:             map[string]store.CachedInviteData{"abc": inv(1, 5, botUser())},
			wantCode:        "abc",
			wantUpdateCodes: []string{"abc"},
		},
		{
			name:     "single-use bot invite disappeared",
			cached:   map[string]store.CachedInviteData{"abc": inv(0, 1, botUser())},
			api:      map[string]store.CachedInviteData{},
			wantCode: "abc",
		},
		{
			name:     "multi-use bot invite disappeared — no match",
			cached:   map[string]store.CachedInviteData{"abc": inv(0, 5, botUser())},
			api:      map[string]store.CachedInviteData{},
			wantCode: "",
		},
		{
			name:     "non-bot single-use invite disappeared — no match",
			cached:   map[string]store.CachedInviteData{"abc": inv(0, 1, otherUser())},
			api:      map[string]store.CachedInviteData{},
			wantCode: "",
		},
		{
			name:            "use count increased on non-bot invite — no match but update recorded",
			cached:          map[string]store.CachedInviteData{"abc": inv(0, 5, otherUser())},
			api:             map[string]store.CachedInviteData{"abc": inv(1, 5, otherUser())},
			wantCode:        "",
			wantUpdateCodes: []string{"abc"},
		},
		{
			name:     "no change in use count",
			cached:   map[string]store.CachedInviteData{"abc": inv(0, 1, botUser())},
			api:      map[string]store.CachedInviteData{"abc": inv(0, 1, botUser())},
			wantCode: "",
		},
		{
			name:     "nil inviter disappeared — no match",
			cached:   map[string]store.CachedInviteData{"abc": inv(0, 1, nil)},
			api:      map[string]store.CachedInviteData{},
			wantCode: "",
		},
		{
			name:     "empty cache",
			cached:   map[string]store.CachedInviteData{},
			api:      map[string]store.CachedInviteData{"abc": inv(1, 5, otherUser())},
			wantCode: "",
		},
		{
			name:     "empty cache and empty api",
			cached:   map[string]store.CachedInviteData{},
			api:      map[string]store.CachedInviteData{},
			wantCode: "",
		},
		{
			name: "multiple invites — only consumed one matched",
			cached: map[string]store.CachedInviteData{
				"alive": inv(0, 1, botUser()),
				"gone":  inv(0, 1, botUser()),
			},
			api: map[string]store.CachedInviteData{
				"alive": inv(0, 1, botUser()),
			},
			wantCode:         "gone",
			wantNoUpdateCode: []string{"alive", "gone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchUsedInvite(tt.cached, tt.api, botID)

			if result.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", result.Code, tt.wantCode)
			}

			for _, code := range tt.wantUpdateCodes {
				if _, ok := result.StoreUpdates[code]; !ok {
					t.Errorf("expected %q in StoreUpdates", code)
				}
			}

			for _, code := range tt.wantNoUpdateCode {
				if _, ok := result.StoreUpdates[code]; ok {
					t.Errorf("did not expect %q in StoreUpdates", code)
				}
			}
		})
	}
}
