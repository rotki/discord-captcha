package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rotki/discord-captcha/internal/api"
	"github.com/rotki/discord-captcha/internal/bot"
	"github.com/rotki/discord-captcha/internal/config"
	"github.com/rotki/discord-captcha/internal/staticfs"
	"github.com/rotki/discord-captcha/internal/store"
	"github.com/rotki/discord-captcha/internal/version"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "4000"
		}
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", port))
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		os.Exit(0)
	}

	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})))

	slog.Info("starting", "version", version.Version, "git_sha", version.GitSHA)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	var inviteStore store.InviteStore
	if cfg.UseRedis() {
		slog.Info("using Redis store", "host", cfg.RedisHost)
		inviteStore = store.NewRedisStore(cfg.RedisHost, cfg.RedisPassword)
	} else {
		slog.Info("using filesystem store", "dir", "/data/invites")
		fsStore, err := store.NewFSStore("/data/invites")
		if err != nil {
			slog.Error("failed to create FS store", "error", err)
			os.Exit(1)
		}
		inviteStore = fsStore
	}

	stopCleanup := store.StartCleanup(inviteStore)

	discordBot, err := bot.New(cfg, inviteStore)
	if err != nil {
		slog.Error("failed to create bot", "error", err)
		os.Exit(1)
	}

	botErr := make(chan error, 1)
	go func() {
		if err := discordBot.Start(); err != nil {
			botErr <- err
		}
	}()

	staticContent, err := fs.Sub(staticfs.StaticFiles, "files")
	if err != nil {
		slog.Error("failed to create static FS", "error", err)
		os.Exit(1)
	}

	server := api.NewServer(cfg, inviteStore, staticContent)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutting down...")
	case err := <-botErr:
		slog.Error("bot failed to start, shutting down", "error", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
	}

	if err := discordBot.Stop(); err != nil {
		slog.Error("failed to stop bot", "error", err)
	}

	stopCleanup()

	if err := inviteStore.Close(); err != nil {
		slog.Error("failed to close store", "error", err)
	}

	slog.Info("shutdown complete")
}
