package store

import (
	"os"
	"testing"
	"time"
)

func TestFSStoreSetAndIterator(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	invite := CachedInvite{
		Code: "test123",
		Data: CachedInviteData{
			Uses:      0,
			MaxUses:   1,
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
			Inviter:   &CachedUser{ID: "1", Username: "testuser"},
		},
	}

	if err := s.Set(invite); err != nil {
		t.Fatal(err)
	}

	found := false
	for code, data := range s.Iterator() {
		if code == "test123" {
			found = true
			if data.Uses != 0 {
				t.Errorf("expected uses=0, got %d", data.Uses)
			}
			if data.MaxUses != 1 {
				t.Errorf("expected maxUses=1, got %d", data.MaxUses)
			}
			if data.Inviter == nil || data.Inviter.Username != "testuser" {
				t.Error("inviter not preserved")
			}
		}
	}
	if !found {
		t.Error("invite not found in iterator")
	}
}

func TestFSStoreDelete(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	invite := CachedInvite{
		Code: "del123",
		Data: CachedInviteData{
			Uses:      0,
			MaxUses:   1,
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		},
	}

	if err := s.Set(invite); err != nil {
		t.Fatal(err)
	}

	if err := s.Delete("del123"); err != nil {
		t.Fatal(err)
	}

	for code := range s.Iterator() {
		if code == "del123" {
			t.Error("deleted invite still present")
		}
	}
}

func TestFSStoreDeleteNonExistent(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := s.Delete("nonexistent"); err != nil {
		t.Errorf("deleting non-existent should not error, got: %v", err)
	}
}

func TestCleanupExpired(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Write expired invite directly to bypass cleanup-on-set
	invite := CachedInvite{
		Code: "expired1",
		Data: CachedInviteData{
			Uses:      0,
			MaxUses:   1,
			ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		},
	}

	// Write directly to disk
	if err := os.WriteFile(s.filePath(invite.Code), []byte(`{"uses":0,"maxUses":1,"expiresAt":"`+invite.Data.ExpiresAt+`"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	Cleanup(s)

	for code := range s.Iterator() {
		if code == "expired1" {
			t.Error("expired invite should have been cleaned up")
		}
	}
}

func TestCleanupMaxUses(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(s.filePath("maxed"), []byte(`{"uses":1,"maxUses":1,"expiresAt":"never"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	Cleanup(s)

	for code := range s.Iterator() {
		if code == "maxed" {
			t.Error("max-uses invite should have been cleaned up")
		}
	}
}

func TestCleanupNoExpiration(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(s.filePath("noexpiry"), []byte(`{"uses":0,"maxUses":0}`), 0o644); err != nil {
		t.Fatal(err)
	}

	Cleanup(s)

	for code := range s.Iterator() {
		if code == "noexpiry" {
			t.Error("invite without expiration should have been cleaned up")
		}
	}
}

func TestFSStoreRejectsPathTraversal(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	maliciousCodes := []string{
		"../../../etc/passwd",
		"..%2F..%2Fetc%2Fpasswd",
		"code/../../escape",
		"code with spaces",
		"code;rm -rf /",
		"",
	}

	for _, code := range maliciousCodes {
		t.Run("Set_"+code, func(t *testing.T) {
			invite := CachedInvite{
				Code: code,
				Data: CachedInviteData{
					Uses:      0,
					MaxUses:   1,
					ExpiresAt: "never",
				},
			}
			if err := s.Set(invite); err == nil {
				t.Errorf("expected error for malicious code %q", code)
			}
		})

		t.Run("Delete_"+code, func(t *testing.T) {
			if err := s.Delete(code); err == nil {
				t.Errorf("expected error for malicious code %q", code)
			}
		})
	}
}

func TestFSStoreAcceptsValidCodes(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	validCodes := []string{"abc123", "ABC-xyz", "code_with_underscores", "a", "123"}
	for _, code := range validCodes {
		t.Run(code, func(t *testing.T) {
			invite := CachedInvite{
				Code: code,
				Data: CachedInviteData{
					Uses:      0,
					MaxUses:   1,
					ExpiresAt: "never",
				},
			}
			if err := s.Set(invite); err != nil {
				t.Errorf("expected no error for valid code %q, got: %v", code, err)
			}
		})
	}
}

func TestCleanupNeverExpires(t *testing.T) {
	dir := t.TempDir()
	s, err := NewFSStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(s.filePath("permanent"), []byte(`{"uses":0,"maxUses":0,"expiresAt":"never"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	Cleanup(s)

	found := false
	for code := range s.Iterator() {
		if code == "permanent" {
			found = true
		}
	}
	if !found {
		t.Error("permanent invite should not be cleaned up")
	}
}
