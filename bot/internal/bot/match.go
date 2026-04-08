package bot

import "github.com/rotki/discord-captcha/internal/store"

// MatchResult holds the outcome of invite matching.
type MatchResult struct {
	// Code is the invite code that was consumed, empty if none matched.
	Code string
	// StoreUpdates contains invites whose use count increased and need
	// to be written back to the cache.
	StoreUpdates map[string]store.CachedInviteData
}

// matchUsedInvite compares a snapshot of cached invites against a snapshot of
// the current Discord API invites to determine which bot invite was consumed.
//
// Detection rules (evaluated per cached invite):
//  1. Use-count increase: invite exists in both maps and api.Uses > cached.Uses,
//     and the inviter is the bot → match.
//  2. Disappeared single-use: invite exists in cached but NOT in api,
//     cached.MaxUses == 1, and the inviter is the bot → match (consumed & auto-deleted).
func matchUsedInvite(
	cached map[string]store.CachedInviteData,
	api map[string]store.CachedInviteData,
	botUserID string,
) MatchResult {
	result := MatchResult{StoreUpdates: make(map[string]store.CachedInviteData)}

	for code, cachedData := range cached {
		apiData, inAPI := api[code]

		if inAPI && apiData.Uses > cachedData.Uses {
			result.StoreUpdates[code] = apiData

			if result.Code == "" && cachedData.Inviter != nil && cachedData.Inviter.ID == botUserID {
				result.Code = code
			}
			continue
		}

		if !inAPI && cachedData.MaxUses == 1 && cachedData.Inviter != nil && cachedData.Inviter.ID == botUserID {
			if result.Code == "" {
				result.Code = code
			}
		}
	}

	return result
}
