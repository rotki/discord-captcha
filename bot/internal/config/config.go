package config

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DiscordToken     string
	DiscordAppID     string
	DiscordGuildID   string
	DiscordChannelID string
	DiscordRoleID    string
	RecaptchaSiteKey string
	RecaptchaSecret  string
	RedisHost        string
	RedisPassword    string
	SiteURL          string
	Port             string
}

func Load() (*Config, error) {
	cfg := &Config{
		DiscordAppID:     os.Getenv("DISCORD_APP_ID"),
		DiscordGuildID:   os.Getenv("DISCORD_GUILD_ID"),
		DiscordChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
		DiscordRoleID:    os.Getenv("DISCORD_ROLE_ID"),
		RecaptchaSiteKey: os.Getenv("RECAPTCHA_SITE_KEY"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		SiteURL:          os.Getenv("SITE_URL"),
		Port:             os.Getenv("PORT"),
	}

	var err error
	if cfg.DiscordToken, err = secret("DISCORD_TOKEN"); err != nil {
		return nil, err
	}
	if cfg.RecaptchaSecret, err = secret("RECAPTCHA_SECRET"); err != nil {
		return nil, err
	}
	if cfg.RedisPassword, err = secret("REDIS_PASSWORD"); err != nil {
		return nil, err
	}

	if cfg.Port == "" {
		cfg.Port = "4000"
	}

	if port, err := strconv.Atoi(cfg.Port); err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("PORT must be a valid port number (1-65535), got %q", cfg.Port)
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN (or DISCORD_TOKEN_FILE) is required")
	}
	if cfg.DiscordAppID == "" {
		return nil, fmt.Errorf("DISCORD_APP_ID is required")
	}
	if cfg.DiscordGuildID == "" {
		return nil, fmt.Errorf("DISCORD_GUILD_ID is required")
	}
	if cfg.DiscordChannelID == "" {
		return nil, fmt.Errorf("DISCORD_CHANNEL_ID is required")
	}
	if cfg.DiscordRoleID == "" {
		return nil, fmt.Errorf("DISCORD_ROLE_ID is required")
	}
	if cfg.RecaptchaSecret == "" {
		return nil, fmt.Errorf("RECAPTCHA_SECRET (or RECAPTCHA_SECRET_FILE) is required")
	}

	return cfg, nil
}

func (c *Config) UseRedis() bool {
	return c.RedisHost != ""
}

// maxSecretSize bounds what is read from a secret file; the secrets here are
// tokens and passwords of well under a kilobyte.
const maxSecretSize = 64 << 10

// secret reads a secret from NAME_FILE (a mounted secret such as
// /run/secrets/name, the Docker convention) or from NAME. Setting both is an
// error rather than a silent choice. A file's trailing newline is dropped, and
// a file with nothing else is an error: an optional secret such as
// REDIS_PASSWORD would otherwise pass as unset. Errors never include the
// secret's value.
func secret(name string) (string, error) {
	value, path := os.Getenv(name), os.Getenv(name+"_FILE")
	if path == "" {
		return value, nil
	}
	if value != "" {
		return "", fmt.Errorf("set either %s or %s_FILE, not both", name, name)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("reading %s_FILE: %w", name, err)
	}
	defer func() { _ = f.Close() }()
	data, err := io.ReadAll(io.LimitReader(f, maxSecretSize+1))
	if err != nil {
		return "", fmt.Errorf("reading %s_FILE: %w", name, err)
	}
	if len(data) > maxSecretSize {
		return "", fmt.Errorf("%s_FILE is larger than %d bytes", name, maxSecretSize)
	}
	value = strings.TrimRight(string(data), "\r\n")
	if value == "" {
		return "", fmt.Errorf("%s_FILE is empty", name)
	}
	return value, nil
}
