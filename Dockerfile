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
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd

# --- Final ---
FROM scratch
COPY --from=go-build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-build /build/server /server
EXPOSE 4000
VOLUME /data
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["/server", "healthcheck"]
ENTRYPOINT ["/server"]
