package httpapi

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type middleware func(http.Handler) http.Handler

func chain(h http.Handler, mws ...middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func withRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				writeError(w, http.StatusInternalServerError, "internal", "Something went wrong.", "")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func withRequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.status, time.Since(start).Round(time.Millisecond))
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// CORSConfig controls cross-origin access for local Hugo development.
type CORSConfig struct {
	Origin string
}

func withCORS(cfg CORSConfig) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", cfg.Origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Vary", "Origin")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// simpleIPRateLimiter is a tiny per-IP token bucket for demo abuse protection.
type simpleIPRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]time.Time
	minGap   time.Duration
}

func newIPRateLimiter(minGap time.Duration) *simpleIPRateLimiter {
	return &simpleIPRateLimiter{
		visitors: make(map[string]time.Time),
		minGap:   minGap,
	}
}

func (l *simpleIPRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if last, ok := l.visitors[ip]; ok && now.Sub(last) < l.minGap {
		return false
	}
	l.visitors[ip] = now
	return true
}

func withRateLimit(limiter *simpleIPRateLimiter) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v1/audits") {
				ip := r.RemoteAddr
				if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
					ip = strings.TrimSpace(strings.Split(fwd, ",")[0])
				}
				if !limiter.allow(ip) {
					writeError(w, http.StatusTooManyRequests, "rate_limited", "Too many requests. Try again shortly.", "")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func withGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Keep middleware present for Phase 6; stdlib gzip wrapper applied selectively in static serving.
		next.ServeHTTP(w, r)
	})
}
