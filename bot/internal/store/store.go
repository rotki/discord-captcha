package store

import (
	"io"
	"iter"
	"log/slog"
	"time"
)

type CachedUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type CachedInviteData struct {
	Uses      int         `json:"uses"`
	MaxUses   int         `json:"maxUses"`
	Inviter   *CachedUser `json:"inviter,omitempty"`
	ExpiresAt string      `json:"expiresAt,omitempty"`
}

type CachedInvite struct {
	Code string           `json:"code"`
	Data CachedInviteData `json:"data"`
}

type InviteStore interface {
	io.Closer
	Set(invite CachedInvite) error
	Delete(code string) error
	Iterator() iter.Seq2[string, CachedInviteData]
}

const cleanupInterval = 5 * time.Minute

// StartCleanup runs periodic cleanup in a goroutine and returns a stop function.
func StartCleanup(s InviteStore) func() {
	ticker := time.NewTicker(cleanupInterval)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				Cleanup(s)
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	return func() { close(done) }
}

func Cleanup(s InviteStore) {
	now := time.Now()
	var toDelete []string

	for code, data := range s.Iterator() {
		switch data.ExpiresAt {
		case "never":
			if data.MaxUses > 0 && data.Uses >= data.MaxUses {
				slog.Info("invite reached max uses, purging", "code", code)
				toDelete = append(toDelete, code)
			}
		case "":
			slog.Info("invite didn't have expiration data, purging", "code", code)
			toDelete = append(toDelete, code)
		default:
			expirationDate, err := time.Parse(time.RFC3339, data.ExpiresAt)
			if err != nil {
				slog.Info("invite had invalid expiration, purging", "code", code, "expiresAt", data.ExpiresAt)
				toDelete = append(toDelete, code)
				continue
			}
			if expirationDate.Before(now) {
				slog.Info("invite expired, purging", "code", code, "expiresAt", data.ExpiresAt)
				toDelete = append(toDelete, code)
			}
		}
	}

	for _, code := range toDelete {
		if err := s.Delete(code); err != nil {
			slog.Error("failed to delete invite during cleanup", "code", code, "error", err)
		}
	}
}
