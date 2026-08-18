package httpapi

import (
	"net/http"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/session"
)

// Options configures the HTTP router.
type Options struct {
	Service    *audit.Service
	Sessions   *session.Manager
	CORSOrigin string
	StaticDir  string
}

// NewRouter builds the application mux with middleware.
func NewRouter(opts Options) http.Handler {
	srv := NewServer(opts.Service, opts.Sessions)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/healthz", srv.handleHealthz)
	mux.HandleFunc("POST /api/v1/audits", srv.handleCreateAudit)
	mux.HandleFunc("GET /api/v1/audits", srv.handleListAudits)
	mux.HandleFunc("GET /api/v1/audits/{id}", srv.handleGetAudit)
	mux.HandleFunc("POST /api/v1/session", srv.handleLogin)
	mux.HandleFunc("DELETE /api/v1/session", srv.handleLogout)
	mux.HandleFunc("GET /api/v1/session", srv.handleSession)

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
