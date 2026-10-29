package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var secretNames = []string{"DISCORD_TOKEN", "RECAPTCHA_SECRET", "REDIS_PASSWORD"}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DISCORD_TOKEN", "test-token")
	t.Setenv("DISCORD_APP_ID", "test-app-id")
	t.Setenv("DISCORD_GUILD_ID", "test-guild-id")
	t.Setenv("DISCORD_CHANNEL_ID", "test-channel-id")
	t.Setenv("DISCORD_ROLE_ID", "test-role-id")
	t.Setenv("RECAPTCHA_SECRET", "test-secret")
	t.Setenv("REDIS_PASSWORD", "")
	for _, name := range secretNames {
		t.Setenv(name+"_FILE", "")
	}
}

func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadFromSecretFiles(t *testing.T) {
	setRequiredEnv(t)
	files := map[string]string{
		"DISCORD_TOKEN":    "file-token\n",
		"RECAPTCHA_SECRET": "file-secret\r\n",
		"REDIS_PASSWORD":   "file-password",
	}
	for name, content := range files {
		t.Setenv(name, "")
		t.Setenv(name+"_FILE", writeFile(t, strings.ToLower(name), content))
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.DiscordToken != "file-token" || cfg.RecaptchaSecret != "file-secret" || cfg.RedisPassword != "file-password" {
		t.Errorf("unexpected secrets: %q %q %q", cfg.DiscordToken, cfg.RecaptchaSecret, cfg.RedisPassword)
	}
}

func TestSecretRejectsBothSources(t *testing.T) {
	for _, name := range secretNames {
		t.Run(name, func(t *testing.T) {
			setRequiredEnv(t)
			t.Setenv(name, "env-value")
			t.Setenv(name+"_FILE", writeFile(t, "s", "file-value"))

			_, err := Load()
			if err == nil || !strings.Contains(err.Error(), "not both") {
				t.Fatalf("expected a not both error, got: %v", err)
			}
			if strings.Contains(err.Error(), "env-value") || strings.Contains(err.Error(), "file-value") {
				t.Errorf("error leaks a secret: %v", err)
			}
		})
	}
}

func TestSecretFileErrors(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DISCORD_TOKEN", "")
	t.Setenv("DISCORD_TOKEN_FILE", filepath.Join(t.TempDir(), "missing"))
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "DISCORD_TOKEN_FILE") {
		t.Errorf("missing file: %v", err)
	}

	t.Setenv("DISCORD_TOKEN_FILE", writeFile(t, "big", strings.Repeat("x", maxSecretSize+1)))
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "larger than") {
		t.Errorf("oversized file: %v", err)
	}

	t.Setenv("DISCORD_TOKEN_FILE", writeFile(t, "empty", "\n"))
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "DISCORD_TOKEN_FILE is empty") {
		t.Errorf("empty file: %v", err)
	}
}

// An optional secret has no required check after it, so only the reader can
// catch an empty file.
func TestOptionalSecretEmptyFile(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("REDIS_PASSWORD_FILE", writeFile(t, "empty", ""))
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "REDIS_PASSWORD_FILE is empty") {
		t.Errorf("empty file: %v", err)
	}
}

func TestLoadSuccess(t *testing.T) {
	setRequiredEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.DiscordToken != "test-token" {
		t.Errorf("expected DiscordToken=test-token, got %q", cfg.DiscordToken)
	}
	if cfg.Port != "4000" {
		t.Errorf("expected default port 4000, got %q", cfg.Port)
	}
}

func TestLoadCustomPort(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("PORT", "8080")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("expected port 8080, got %q", cfg.Port)
	}
}

func TestLoadMissingRequired(t *testing.T) {
	required := []string{
		"DISCORD_TOKEN",
		"DISCORD_APP_ID",
		"DISCORD_GUILD_ID",
		"DISCORD_CHANNEL_ID",
		"DISCORD_ROLE_ID",
		"RECAPTCHA_SECRET",
	}

	for _, env := range required {
		t.Run(env, func(t *testing.T) {
			setRequiredEnv(t)
			t.Setenv(env, "")

			_, err := Load()
			if err == nil {
				t.Errorf("expected error when %s is missing", env)
			}
		})
	}
}

func TestLoadInvalidPort(t *testing.T) {
	invalidPorts := []string{"abc", "0", "65536", "-1", "99999"}
	for _, port := range invalidPorts {
		t.Run(port, func(t *testing.T) {
			setRequiredEnv(t)
			t.Setenv("PORT", port)

			_, err := Load()
			if err == nil {
				t.Errorf("expected error for invalid port %q", port)
			}
		})
	}
}

func TestUseRedis(t *testing.T) {
	cfg := &Config{RedisHost: ""}
	if cfg.UseRedis() {
		t.Error("expected UseRedis=false when RedisHost is empty")
	}

	cfg.RedisHost = "localhost"
	if !cfg.UseRedis() {
		t.Error("expected UseRedis=true when RedisHost is set")
	}
}
