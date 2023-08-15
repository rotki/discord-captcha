FROM node:18 as builder

COPY ./ /build/

WORKDIR /build

RUN --mount=type=cache,target=/root/.npm/_cacache/ \
    --mount=type=cache,target=/root/.local/share/pnpm/store \
    npm install -g pnpm@8 && \
    pnpm install --frozen-lockfile && \
    pnpm run build

FROM node:18-alpine as production

WORKDIR /app
ARG GIT_SHA
ENV GIT_SHA=${GIT_SHA}
ENV NITRO_HOST=0.0.0.0
ENV NITRO_PORT=4000

RUN --mount=type=cache,target=/root/.npm/_cacache/ \
    npm install -g pm2@5.3.0 && \
    mkdir data

COPY --from=builder /build/.output ./.output/
COPY --from=builder /build/ecosystem.config.cjs ./
COPY --from=builder /build/package.json ./

EXPOSE ${NITRO_PORT}
VOLUME /app/data

CMD ["pm2-runtime", "ecosystem.config.cjs"]
HEALTHCHECK --start-period=30s --retries=2 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${NITRO_PORT}/health || exit 1
