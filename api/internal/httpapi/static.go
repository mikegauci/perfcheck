package httpapi

import (
	"net/http"
	"path"
	"strings"
)

func staticHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		if isFingerprintedAsset(r.URL.Path) {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		fileServer.ServeHTTP(w, r)
	})
}

func isFingerprintedAsset(p string) bool {
	base := path.Base(p)
	if !(strings.HasSuffix(base, ".css") || strings.HasSuffix(base, ".js")) {
		return false
	}
	// Hugo fingerprints as name.<hash>.ext
	parts := strings.Split(base, ".")
	return len(parts) >= 3
}
