package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/httpapi"
	"github.com/mikegauci/perfcheck/api/internal/storage"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	corsOrigin := flag.String("cors-origin", "http://localhost:1313", "Allowed CORS origin (empty disables CORS)")
	staticDir := flag.String("static", "", "Optional directory of Hugo public/ files to serve")
	flag.Parse()

	repo := storage.NewMemory(100)
	svc := audit.NewService(repo)
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
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("perfcheckd listening on %s (cors=%q static=%q)", *addr, *corsOrigin, *staticDir)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
