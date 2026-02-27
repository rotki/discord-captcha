# --- Build Frontend ---
FROM node:24-alpine AS web-build
WORKDIR /build
COPY web/package.json web/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY web/ .
ARG VITE_RECAPTCHA_SITE_KEY
ARG VITE_SITE_URL
ENV VITE_RECAPTCHA_SITE_KEY=${VITE_RECAPTCHA_SITE_KEY}
ENV VITE_SITE_URL=${VITE_SITE_URL}
RUN pnpm build

# --- Build Go ---
FROM golang:1.26-alpine AS go-build
WORKDIR /build
COPY bot/go.mod bot/go.sum ./
RUN go mod download
COPY bot/ .
COPY --from=web-build /build/dist ./internal/staticfs/files/
ARG GIT_SHA=unknown
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X github.com/rotki/discord-captcha/internal/version.Version=${VERSION} -X github.com/rotki/discord-captcha/internal/version.GitSHA=${GIT_SHA}" \
    -o server ./cmd

# --- Final ---
FROM scratch
ARG GIT_SHA=unknown
ARG VERSION=dev
LABEL org.opencontainers.image.title="discord-captcha"
LABEL org.opencontainers.image.description="Discord bot with captcha verification for server access"
LABEL org.opencontainers.image.source="https://github.com/rotki/discord-captcha"
LABEL org.opencontainers.image.licenses="AGPL-3.0-or-later"
LABEL org.opencontainers.image.revision="${GIT_SHA}"
LABEL org.opencontainers.image.version="${VERSION}"

COPY --from=go-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-build /build/server /server
EXPOSE 4000
VOLUME /data
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/server", "healthcheck"]
ENTRYPOINT ["/server"]
