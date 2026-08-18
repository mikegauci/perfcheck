.PHONY: help dev api hugo test test-api test-frontend typecheck build lint clean

help:
	@echo "PerfCheck targets:"
	@echo "  make dev            - run Hugo (1313) and API (8080) together"
	@echo "  make hugo           - run Hugo server only"
	@echo "  make api            - run Go API only"
	@echo "  make test           - run Go and frontend tests"
	@echo "  make test-api       - run Go tests with race detector"
	@echo "  make test-frontend  - run Vitest"
	@echo "  make typecheck      - TypeScript check"
	@echo "  make build          - build Hugo site and Go binary"
	@echo "  make lint           - go vet, typecheck, ESLint, Stylelint"
	@echo "  make clean          - remove build artefacts"

dev:
	@echo "Start API in one terminal: make api"
	@echo "Start Hugo in another:     make hugo"
	@echo "Or run both with:          (make api &) && make hugo"

hugo:
	hugo server --bind 127.0.0.1 --port 1313 --disableFastRender

api:
	cd api && go run ./cmd/perfcheckd

test: test-api test-frontend

test-api:
	cd api && go test ./... -race -count=1

test-frontend:
	npm test

typecheck:
	npm run typecheck

build:
	npm ci
	hugo --minify
	cd api && CGO_ENABLED=0 go build -o bin/perfcheckd ./cmd/perfcheckd

lint:
	cd api && go vet ./...
	npm run typecheck
	npm run lint
	npm run lint:css

clean:
	rm -rf public resources/_gen api/bin coverage
