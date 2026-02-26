package config

import (
	"testing"
)

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DISCORD_TOKEN", "test-token")
	t.Setenv("DISCORD_APP_ID", "test-app-id")
	t.Setenv("DISCORD_GUILD_ID", "test-guild-id")
	t.Setenv("DISCORD_CHANNEL_ID", "test-channel-id")
	t.Setenv("DISCORD_ROLE_ID", "test-role-id")
	t.Setenv("RECAPTCHA_SECRET", "test-secret")
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
