package store

import (
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o644
)

var validInviteCode = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type FSStore struct {
	baseDir string
	mu      sync.RWMutex
}

func NewFSStore(baseDir string) (*FSStore, error) {
	if err := os.MkdirAll(baseDir, dirPermissions); err != nil {
		return nil, err
	}
	return &FSStore{baseDir: baseDir}, nil
}

func (s *FSStore) validateCode(code string) error {
	if !validInviteCode.MatchString(code) {
		return fmt.Errorf("invalid invite code: %q", code)
	}
	return nil
}

func (s *FSStore) filePath(code string) string {
	return filepath.Join(s.baseDir, code+".json")
}

func (s *FSStore) Set(invite CachedInvite) error {
	if err := s.validateCode(invite.Code); err != nil {
		return err
	}

	if err := s.writeFile(invite); err != nil {
		return err
	}

	return nil
}

func (s *FSStore) writeFile(invite CachedInvite) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(invite.Data)
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath(invite.Code), data, filePermissions)
}

func (s *FSStore) Delete(code string) error {
	if err := s.validateCode(code); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.Remove(s.filePath(code))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *FSStore) Close() error {
	return nil
}

func (s *FSStore) Iterator() iter.Seq2[string, CachedInviteData] {
	return func(yield func(string, CachedInviteData) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		entries, err := os.ReadDir(s.baseDir)
		if err != nil {
			slog.Error("failed to read invite directory", "error", err)
			return
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
				continue
			}

			code := strings.TrimSuffix(entry.Name(), ".json")
			data, err := os.ReadFile(filepath.Join(s.baseDir, entry.Name()))
			if err != nil {
				slog.Error("failed to read invite file", "code", code, "error", err)
				continue
			}

			var inviteData CachedInviteData
			if err := json.Unmarshal(data, &inviteData); err != nil {
				slog.Error("failed to parse invite file", "code", code, "error", err)
				continue
			}

			if !yield(code, inviteData) {
				return
			}
		}
	}
}
