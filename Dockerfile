# Multi-stage build: Hugo site + Go API in one image.
# Build:  docker build -t perfcheck .
# Run:    docker run --rm -p 8080:8080 perfcheck

FROM node:20-bookworm AS assets
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY assets ./assets
COPY layouts ./layouts
COPY content ./content
COPY static ./static
COPY hugo.toml postcss.config.js .browserslistrc ./

FROM hugomods/hugo:exts-0.165.0 AS hugo
# Dart Sass is included in many Hugo images; install if missing.
USER root
RUN command -v sass >/dev/null || (curl -fsSL https://github.com/sass/dart-sass/releases/download/1.77.8/dart-sass-1.77.8-linux-x64.tar.gz | tar -xz -C /usr/local/bin --strip-components=1)
WORKDIR /src
COPY --from=assets /src /src
COPY --from=assets /src/node_modules /src/node_modules
RUN hugo --minify

FROM golang:1.22-bookworm AS api
WORKDIR /src
COPY api/go.mod ./
RUN go mod download
COPY api ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /perfcheckd ./cmd/perfcheckd

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=hugo /src/public /app/public
COPY --from=api /perfcheckd /app/perfcheckd
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/perfcheckd", "-addr", ":8080", "-cors-origin", "", "-static", "/app/public"]
