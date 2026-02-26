package config

import (
	"fmt"
	"os"
	"strconv"
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
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		DiscordAppID:     os.Getenv("DISCORD_APP_ID"),
		DiscordGuildID:   os.Getenv("DISCORD_GUILD_ID"),
		DiscordChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
		DiscordRoleID:    os.Getenv("DISCORD_ROLE_ID"),
		RecaptchaSiteKey: os.Getenv("RECAPTCHA_SITE_KEY"),
		RecaptchaSecret:  os.Getenv("RECAPTCHA_SECRET"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		SiteURL:          os.Getenv("SITE_URL"),
		Port:             os.Getenv("PORT"),
	}

	if cfg.Port == "" {
		cfg.Port = "4000"
	}

	if port, err := strconv.Atoi(cfg.Port); err != nil || port < 1 || port > 65535 {
		return nil, fmt.Errorf("PORT must be a valid port number (1-65535), got %q", cfg.Port)
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
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
		return nil, fmt.Errorf("RECAPTCHA_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) UseRedis() bool {
	return c.RedisHost != ""
}
