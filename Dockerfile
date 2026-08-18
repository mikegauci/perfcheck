# Multi-stage build: Hugo site + Go API in one image.
# Build:  docker build -t perfcheck .
# Run:    docker run --rm -p 8080:8080 perfcheck

FROM node:20-bookworm AS deps
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci

FROM hugomods/hugo:exts AS hugo
USER root
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates \
  && curl -fsSL https://github.com/sass/dart-sass/releases/download/1.77.8/dart-sass-1.77.8-linux-x64.tar.gz \
    | tar -xz -C /usr/local/bin --strip-components=1 \
  && rm -rf /var/lib/apt/lists/*
WORKDIR /src
COPY --from=deps /src/node_modules ./node_modules
COPY package.json package-lock.json postcss.config.js .browserslistrc hugo.toml ./
COPY assets ./assets
COPY layouts ./layouts
COPY content ./content
COPY static ./static
RUN hugo --minify

FROM golang:1.22-bookworm AS api
WORKDIR /src
COPY api/go.mod ./
RUN go mod download 2>/dev/null || true
COPY api ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /perfcheckd ./cmd/perfcheckd

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=hugo /src/public /app/public
COPY --from=api /perfcheckd /app/perfcheckd
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/perfcheckd", "-addr", ":8080", "-cors-origin", "", "-static", "/app/public"]
