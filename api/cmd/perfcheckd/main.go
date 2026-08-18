package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/httpapi"
	"github.com/mikegauci/perfcheck/api/internal/session"
	"github.com/mikegauci/perfcheck/api/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	corsOrigin := flag.String("cors-origin", "http://localhost:1313", "Allowed CORS origin (empty disables CORS)")
	staticDir := flag.String("static", "", "Optional directory of Hugo public/ files to serve")
	cacheTTL := flag.Duration("cache-ttl", 15*time.Minute, "TTL for scoring cache")
	dbPath := flag.String("db", "", "SQLite path (empty uses in-memory storage)")
	dashboardPassword := flag.String("dashboard-password", "", "Optional password protecting GET /api/v1/audits")
	sessionSecret := flag.String("session-secret", os.Getenv("PERFCHECK_SESSION_SECRET"), "HMAC secret for session cookies")
	flag.Parse()

	scorer := audit.NewConfiguredScorer(audit.ScorerOptions{
		CacheTTL: *cacheTTL,
	})

	repo, err := storage.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if c, ok := repo.(interface{ Close() error }); ok {
		defer c.Close()
	}

	sessions := session.New(*sessionSecret, *dashboardPassword, false)
	svc := audit.NewService(repo, scorer)
	handler := httpapi.NewRouter(httpapi.Options{
		Service:    svc,
		Sessions:   sessions,
		CORSOrigin: *corsOrigin,
		StaticDir:  *staticDir,
	})

	server := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("perfcheckd listening on %s (cors=%q static=%q db=%q)", *addr, *corsOrigin, *staticDir, *dbPath)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
