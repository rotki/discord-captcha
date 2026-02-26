# @rotki/discord-captcha

A Discord captcha verification bot and website. Users solve a reCAPTCHA to receive a single-use Discord invite link. The bot monitors invite usage and assigns a role to members who join through captcha-verified invites.

## Architecture

- **`bot/`** — Go backend: Discord bot (gateway + slash commands) and HTTP server (reCAPTCHA verification, invite creation, static file serving)
- **`web/`** — Vue 3 + Vite frontend: SPA with `@rotki/ui-library`, Tailwind CSS, and `vue-i18n`
- Static frontend is embedded into the Go binary at build time via `//go:embed`
- Deployed as a single binary in a `scratch` Docker container

## Configuration

### Create and configure

Create a new [Discord application](https://discord.com/developers/applications).

In the bot settings of your new application make sure to go under `Privileged Gateway Intents`
and enable `SERVER_MEMBERS_INTENT`.

### Permissions

Make sure the bot has the following permissions:

- Manage Server
- Manage Roles
- Manage Channels
- Create Instant Invite

### Add the bot to your server

Make sure that your bot is not public.
To invite the bot to your server you can use the following link after replacing the placeholders:

```
https://discord.com/oauth2/authorize?client_id={appId}&permissions={permissions}&scope=bot
```

For the permissions use the permission calculator in the application bot settings.

### Role Addition

You should ensure that the added role is below the bot's role in the Guild's
roles settings.

Otherwise, this might cause permission errors even if you have given the `MANAGE_ROLES`
permission to the bot.

## Environment Variables

Copy `.env.example` to `.env` and fill in the values:

| Variable | Required | Description |
|---|---|---|
| `DISCORD_TOKEN` | Yes | Discord bot token |
| `DISCORD_APP_ID` | Yes | Discord application ID |
| `DISCORD_GUILD_ID` | Yes | Target Discord server ID |
| `DISCORD_CHANNEL_ID` | Yes | Channel to create invites in |
| `DISCORD_ROLE_ID` | Yes | Role to assign on verified join |
| `RECAPTCHA_SECRET` | Yes | Google reCAPTCHA v2 secret key |
| `VITE_RECAPTCHA_SITE_KEY` | Yes (build-time) | Google reCAPTCHA v2 site key |
| `VITE_SITE_URL` | No (build-time) | Main site URL (default: `https://rotki.com`) |
| `PORT` | No | HTTP server port (default: `4000`, must be 1-65535) |
| `LOG_LEVEL` | No | Log level: `INFO` (default) or `DEBUG` |
| `REDIS_HOST` | No | Redis host for invite cache |
| `REDIS_PASSWORD` | No | Redis password |

`VITE_`-prefixed variables are only needed at frontend build time.
If `REDIS_HOST` is not set, the filesystem store (`/data/invites`) is used.

## Development

### Prerequisites

- Go 1.26+
- Node.js 24+ with pnpm 10+
- GNU Make (optional, for convenience targets)

### Setup

```bash
cp .env.example .env
# Fill in your Discord and reCAPTCHA credentials in .env
```

### Using Make

The root `Makefile` automatically loads `.env` and provides convenience targets:

```bash
make dev          # Start bot and web dev server in parallel
make dev-bot      # Build and run the Go bot
make dev-web      # Start the Vite dev server
make build        # Build both bot and web
make test         # Run all tests (Go + web typecheck)
make lint         # Lint both bot and web
make clean        # Clean build artifacts
make docker       # Build docker image
make docker-up    # Start with docker compose
make docker-down  # Stop docker compose
make help         # Show all available targets
```

The `bot/` directory has its own `Makefile` with Go-specific targets:

```bash
cd bot
make build        # Build the Go binary
make test         # Run tests
make test-cover   # Run tests with coverage
make lint         # Run golangci-lint
make fmt          # Format code
make vet          # Run go vet
make tidy         # Tidy go.mod
```

### Running locally (dev mode)

The easiest way is via Make:

```bash
make dev
```

Or run the Go backend and Vite dev server in separate terminals:

**Terminal 1 — Go backend:**

```bash
cd bot
set -a && source ../.env && set +a
go run ./cmd
```

**Terminal 2 — Frontend dev server:**

```bash
cd web
pnpm install
set -a && source ../.env && set +a
pnpm dev
```

Visit `http://localhost:5173`. The page renders with hot-reload, and captcha
submissions are proxied to the Go backend at `:4000`.

### Running locally (production-like)

This builds the frontend and embeds it into the Go binary, matching how Docker
runs it.

```bash
# 1. Build the frontend
cd web
pnpm install
set -a && source ../.env && set +a
pnpm build

# 2. Copy build output into the Go embed directory
rm -f ../bot/internal/staticfs/files/placeholder
cp -r dist/* ../bot/internal/staticfs/files/

# 3. Build and run the Go server
cd ../bot
set -a && source ../.env && set +a
go run ./cmd
```

Visit `http://localhost:4000`.

### Frontend commands

```bash
cd web
pnpm dev         # dev server on http://localhost:5173
pnpm build       # production build to web/dist/
pnpm lint        # eslint
pnpm typecheck   # vue-tsc
```

### Verifying the setup

```bash
# Health check
curl http://localhost:4000/health
# Should return: Ok

# Test captcha flow
# Visit the site, complete the reCAPTCHA, and verify an invite link is generated
```

## Production

### Docker (recommended)

```bash
docker compose up --build
```

This builds a multi-stage image (Node for frontend, Go for backend) producing a
minimal `scratch` container. The app is available at `http://localhost:4000`.

Build-time variables (`VITE_RECAPTCHA_SITE_KEY`, `VITE_SITE_URL`) are passed as
Docker build args. Runtime variables are set via `environment` in
`docker-compose.yml`.

To enable debug logging:

```bash
LOG_LEVEL=DEBUG docker compose up
```

### Manual

```bash
# 1. Build frontend
cd web
VITE_RECAPTCHA_SITE_KEY=... VITE_SITE_URL=https://rotki.com pnpm build

# 2. Copy frontend into Go embed directory
rm -f ../bot/internal/staticfs/files/placeholder
cp -r dist/* ../bot/internal/staticfs/files/

# 3. Build Go binary
cd ../bot
CGO_ENABLED=0 go build -ldflags="-s -w" -o server ./cmd

# 4. Run (with env vars set)
./server
```

## Slash Commands

| Command | Description |
|---|---|
| `/logsdir` | Links to the log files location documentation |
| `/datadir` | Links to the data directory location documentation |
