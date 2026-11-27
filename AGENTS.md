# rotki/discord-captcha

A Discord captcha verification bot and website. Users solve a reCAPTCHA on the site to receive a
single-use invite link, and the bot assigns a role to members who join through such an invite. It
ships as one Go binary with the Vue frontend embedded, in a `scratch` Docker image.

This file is the maintenance guide for humans and coding agents. `README.md` covers configuration,
environment variables, and running the project locally and in Docker; this guide links there
instead of repeating it.

## Toolchain

- **Go**, the version in `bot/go.mod` (`go` directive). The `golang:<version>-alpine` image in the
  `Dockerfile` names the same version.
- **Node.js 24** (`.nvmrc`). The `node:<version>-alpine` image in the `Dockerfile` names the same
  version.
- **pnpm 12**, pinned through `packageManager` in `web/package.json`. There is no root
  `package.json`; every pnpm command runs in `web/`.
- **golangci-lint**, the version in `.github/workflows/ci.yml` (the `# renovate:` comment above it
  keeps it updated). Install the same version locally.
- **Docker** with Compose, to build the production image.

## Layout

- `bot/cmd/`: the server entry point. `/server healthcheck` is the container health check.
- `bot/internal/api/`: HTTP server: reCAPTCHA verification, invite creation, static file serving.
- `bot/internal/bot/`: Discord gateway handling, slash commands, invite matching and role
  assignment.
- `bot/internal/config/`: environment configuration, including the `<NAME>_FILE` secret variants.
- `bot/internal/store/`: invite store: Redis when `REDIS_HOST` is set, the filesystem otherwise.
- `bot/internal/staticfs/`: the `//go:embed` of the built frontend. `files/placeholder` keeps the
  directory in git so the package builds without a frontend build; the Docker build copies
  `web/dist/` over it.
- `bot/internal/version/`: version and git SHA, set through `-ldflags` at build time.
- `web/`: Vue 3, Vite, `@rotki/ui-library`, Tailwind CSS and `vue-i18n` (`web/src/locales/`). It
  also holds the commitlint config.
- `Dockerfile`: three stages: frontend build, Go build, `scratch` runtime as user `65532`.
- `Makefile` (root) and `bot/Makefile`: convenience targets. The root one loads `.env` if present.

## Commands

| Command | What it does |
| --- | --- |
| `make build` | Build the Go binary and the frontend |
| `make test` | Go tests, then the frontend typecheck |
| `make lint` | golangci-lint, then ESLint on `web/` |
| `make -C bot vet` / `make -C bot fmt` | `go vet` / `gofmt -w` |
| `make -C bot test-cover` | Go tests with a coverage summary |
| `cd web && pnpm run typecheck` | `vue-tsc` |
| `cd web && pnpm run build` | Production frontend build into `web/dist/` |
| `make docker` | Build the image with Compose |

CI (`.github/workflows/ci.yml`) runs commitlint, then the `go` job (golangci-lint, tests, build),
the `web` job (lint, typecheck, build) and `build-docker` (image build without pushing).

## Conventions

- Commits follow Conventional Commits, with `bot` or `web` as the scope where it fits
  (`feat(bot): ...`, `fix(web): ...`). commitlint runs in CI only; there are no local git hooks,
  so check messages with `cd web && pnpm exec commitlint --from HEAD~1`. `feat` and `fix` commits
  end up in the generated changelog.
- JavaScript dependencies are pinned exactly (`saveExact: true` in `web/pnpm-workspace.yaml`).
- A new environment variable is added in `bot/internal/config/`, its test, `.env.example` and the
  table in `README.md` together. Secrets also get a `<NAME>_FILE` variant like the existing ones.

## Hard rules

Follow these in every change.

- **Never read, print or commit real secrets.** `.env` holds live Discord and reCAPTCHA credentials
  and stays out of git and out of tool output. Work from `.env.example` and the README table.
- **Keep `bot/internal/staticfs/files/placeholder` committed, and never commit a built frontend
  there.** The Docker build fills that directory; a committed build would go stale.
- **Keep version pins in step.** The Go version in `bot/go.mod` and the `golang` image in the
  `Dockerfile`, and `.nvmrc` and the `node` image, change together.
- **Do not widen the bot's Discord permissions or intents quietly.** The README lists what the bot
  needs; a change there is called out in the pull request.
- **Do not lower a check to make it pass.** No `//nolint`, `eslint-disable`, `@ts-expect-error`,
  skipped tests or zizmor ignores, unless the commit explains why the check is wrong for this case.

## Verify before you finish

A change is done when the checks for it pass locally. Report which commands you ran and their
result; if you skipped one, say so and why.

| You changed | Run |
| --- | --- |
| Anything in `bot/` | `make -C bot test`, `make -C bot lint`, `make -C bot vet`; a test that fails without your change |
| Anything in `web/` | `cd web && pnpm run lint && pnpm run typecheck && pnpm run build` |
| Behavior a user sees | Run it (see "Running locally" in the README) and check the page and the bot |
| `Dockerfile`, dependencies or the embed | `make docker` (or `docker build .`) |
| Configuration | `bot/internal/config/` tests, `.env.example` and the README table |
| `.github/workflows/` | `actionlint` and `GH_TOKEN=$(gh auth token) zizmor .github/workflows`; pin every action to a commit SHA |

When a check fails, find the cause. Do not retry until it passes, and do not change a test to match
broken behavior.

## Git and pull requests

On top of the Conventions above, for agents:

- Branch names: `type/short-description` (`fix/invite-expiry`), never an agent-tool prefix such as
  `claude/`. Pull requests target `main`.
- End every commit you create or amend with exactly one of each trailer, after a blank line:

  ```text
  LLM-Model: <provider>/<model-id>
  Harness: <harness-name> <version>
  ```

  Use the model selected in the current session, and the version of the harness process that runs
  you, not of a CLI on `PATH` (a desktop app ships its own copy). Walk up the parent processes
  (`ps -o pid=,ppid=,args= -p <pid>`) to the binary and ask it for `--version`. If you cannot
  verify either value, ask before committing.
- No AI `Co-Authored-By` trailers and no "generated with" lines in commits, pull requests or code.
- Do not push to `main`, merge pull requests or create tags unless the maintainer asks for that
  specific step.

## Recording what you learn

Before you finish, ask whether you hit something the next person would trip over. If you did, first
try to make a tool catch it: a linter setting, a test or a CI check. Write it down only when no tool
can, short and with the fix or the reason, in this file and in the same pull request. Unfinished
work goes to the maintainer as a proposed GitHub issue, never as a to-do list in this file.

## Dependency updates

Renovate lists updates on the Dependency Dashboard issue and opens pull requests only on approval
(`.github/renovate.json5`). It covers Go modules, the `web/` packages, the Docker base images,
GitHub Actions, and tool versions in workflows (golangci-lint through its `# renovate:` comment,
and `npx <pkg>@<version>` calls).

**Prefer one batched update over approving Renovate pull requests one at a time.** Use the
dashboard as the checklist, and do a maintenance pass by hand when updates have piled up and before
a release. A batch is a single branch and pull request, so CI, the image build and review run once.
Renovate notices the merged versions and clears those dashboard entries by itself. Tick a single
dashboard box only for an update that cannot wait, such as a fix for a security advisory.

### Policy

- **7-day release-age cooldown** in `.github/renovate.json5` and in `web/pnpm-workspace.yaml`. The
  `@rotki/*` packages are excluded.
- **Node types follow the runtime:** Renovate restricts `@types/node` and `node` to 24.x.
- **Overrides:** `web/pnpm-workspace.yaml` overrides `tinypool` under `oxfmt` for security
  advisories; it comes in through `@rotki/eslint-config`. Check on each pass whether the parent has
  caught up, and remove the override once it has.
- GitHub Actions are pinned by commit SHA with the version in a trailing comment.

### Procedure

1. Start from an up to date `main` on a branch such as `chore/deps-update`.
2. Gather the state: the dashboard
   (`gh issue list -R rotki/discord-captcha --search "Dependency Dashboard in:title"`),
   `cd web && pnpm audit && pnpm outdated`, and `cd bot && go list -m -u all`.
3. Check release dates and skip anything younger than 7 days: `npm view <package> "time[<version>]"`
   for npm packages, `go list -m -json <module>@<version>` (the `Time` field) for Go modules.
4. Update:
   - `web/`: `pnpm add -E <package>@<version>`, then `pnpm update --depth Infinity <package>` for
     vulnerable transitive packages.
   - `bot/`: `go get <module>@<version>`, then `go mod tidy`.
   - The Go and Node versions together with their `Dockerfile` images (see Hard rules), pnpm with
     `corepack use pnpm@<version>` in `web/`, golangci-lint in `ci.yml`, and action pins with
     `gh api repos/<owner>/<repo>/commits/<tag> -q .sha`.
5. Verify with the table above, including `make docker`.
6. Commit in parts: `ci: update github actions`, `build(web): ...` or `chore(deps): ...` for the
   web packages, and `build(bot): ...` or `chore(deps): ...` for Go modules, each with its lockfile
   or `go.sum`.

## Workflow changes

Check any change under `.github/workflows/`, including pin bumps, with both tools:

- [actionlint](https://github.com/rhysd/actionlint) checks that workflows are valid: syntax,
  expressions, references, action inputs, and `run:` scripts through shellcheck. Run `actionlint`
  from the repository root.
- [zizmor](https://docs.zizmor.sh/) checks workflow security. Run it with a token so the online
  audits (impostor commits, pin and comment mismatches, known-vulnerable actions) run too:

  ```sh
  GH_TOKEN=$(gh auth token) zizmor .github/workflows
  ```

The current baseline for both is no findings, and a change must keep it that way. Keep the
existing patterns: `permissions: {}` at the top level with per-job grants,
`persist-credentials: false` on checkout, SHA pins with a version comment, and expressions that
carry outside input passed to `run:` through `env:`.

## Releasing

Releases use calendar versions, `vYYYY.M.N`: year, month without a leading zero, and a counter that
starts at 0 each month (`v2026.10.0`, then `v2026.10.1`). There is no version file and no release
commit: the tag goes on the `main` commit being released.

1. Make sure `main` is green and contains everything for the release.
2. Tag the release commit and push the tag to `rotki/discord-captcha`:

   ```sh
   git tag vYYYY.M.N <commit>
   git push <rotki-remote> vYYYY.M.N
   ```

3. `.github/workflows/release.yml` creates the GitHub release with notes from `changelogithub`.
4. Verify with `gh release view vYYYY.M.N -R rotki/discord-captcha`.

### Deploying

CI builds the image on every push but does not publish it. Deployments run from rotki's
infrastructure, outside this repository, and only a maintainer with access to it runs them. Agents
do not deploy.

1. Release first (above), so the deployment builds a tagged commit on `main`.
2. Ask a maintainer to deploy the tag to production. There is no staging deployment of the bot.
   The deployment builds the image from that tag on the server and updates the running bot.
3. Check the deployed bot: it is online in Discord, the captcha page loads, and a new member can
   complete verification.

Production passes the bot's credentials as Docker secrets, mounted under `/run/secrets/` and read
through the `<NAME>_FILE` variables. A new secret therefore needs its `_FILE` variant here before it
can be deployed, and a build older than v2026.10.1, which added them, starts without credentials.

## Repository settings

These live in GitHub, not in the repository. Keep this list current when they change.

- `main` is protected: 1 approving review with code owner review required, and the `commit-lint`,
  `license/cla`, `build-docker`, `go` and `web` checks must pass.
- CI also runs on pushes and pull requests for `staging` and `develop`.
- `CODEOWNERS` assigns everything to `@rotki/maintainers`, which has write access.
- There are no environments or repository secrets; the workflows use only the built-in
  `GITHUB_TOKEN`.
