package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/httpapi"
	"github.com/mikegauci/perfcheck/api/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	corsOrigin := flag.String("cors-origin", "http://localhost:1313", "Allowed CORS origin (empty disables CORS)")
	staticDir := flag.String("static", "", "Optional directory of Hugo public/ files to serve")
	scorerMode := flag.String("scorer", "auto", "Scoring engine: auto, psi, fetch or mock")
	cacheTTL := flag.Duration("cache-ttl", 15*time.Minute, "TTL for scoring cache")
	dbPath := flag.String("db", "", "SQLite path (empty uses in-memory storage)")
	flag.Parse()

	scorer, err := audit.NewConfiguredScorer(audit.ScorerOptions{
		Mode:     *scorerMode,
		PSIKey:   os.Getenv("PSI_API_KEY"),
		CacheTTL: *cacheTTL,
	})
	if err != nil {
		log.Fatal(err)
	}

	repo, err := storage.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if c, ok := repo.(interface{ Close() error }); ok {
		defer c.Close()
	}

	svc := audit.NewService(repo, scorer)
	handler := httpapi.NewRouter(httpapi.Options{
		Service:    svc,
		CORSOrigin: *corsOrigin,
		StaticDir:  *staticDir,
	})

	server := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("perfcheckd listening on %s (cors=%q static=%q scorer=%s db=%q)", *addr, *corsOrigin, *staticDir, *scorerMode, *dbPath)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
