# Multi-stage build: Hugo site + Go API in one image.
# Build:  docker build -t perfcheck .
# Run:    docker run --rm -p 8080:8080 -v perfcheck-data:/data perfcheck

FROM node:20-bookworm AS deps
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci

FROM hugomods/hugo:exts AS hugo
WORKDIR /src
COPY --from=deps /src/node_modules ./node_modules
COPY package.json package-lock.json postcss.config.cjs .browserslistrc hugo.toml ./
COPY config ./config
COPY assets ./assets
COPY layouts ./layouts
COPY content ./content
COPY static ./static
ARG HUGO_BASEURL=http://localhost:1313/
ENV HUGO_PARAMS_APIBASEURL=
RUN hugo --minify -b "${HUGO_BASEURL}"

FROM golang:1.25-bookworm AS api
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /perfcheckd ./cmd/perfcheckd

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=hugo /src/public /app/public
COPY --from=api /perfcheckd /app/perfcheckd
EXPOSE 8080
VOLUME ["/data"]
ENTRYPOINT ["/app/perfcheckd", "-addr", ":8080", "-cors-origin", "", "-static", "/app/public", "-db", "/data/perfcheck.db", "-scorer", "auto"]
