package httpapi

import (
	"net/http"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
)

// Options configures the HTTP router.
type Options struct {
	Service    *audit.Service
	CORSOrigin string
	StaticDir  string
}

// NewRouter builds the application mux with middleware.
func NewRouter(opts Options) http.Handler {
	srv := NewServer(opts.Service)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/healthz", srv.handleHealthz)
	mux.HandleFunc("POST /api/v1/audits", srv.handleCreateAudit)
	mux.HandleFunc("GET /api/v1/audits", srv.handleListAudits)
	mux.HandleFunc("GET /api/v1/audits/{id}", srv.handleGetAudit)

	if opts.StaticDir != "" {
		mux.Handle("/", staticHandler(opts.StaticDir))
	}

	limiter := newIPRateLimiter(200 * time.Millisecond)

	return chain(
		mux,
		withRecover,
		withRequestLog,
		withCORS(CORSConfig{Origin: opts.CORSOrigin}),
		withRateLimit(limiter),
		withGzip,
	)
}
