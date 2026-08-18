package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/session"
	"github.com/mikegauci/perfcheck/api/internal/validate"
)

// Server holds HTTP dependencies.
type Server struct {
	svc      *audit.Service
	sessions *session.Manager
}

// NewServer constructs an API server.
func NewServer(svc *audit.Service, sessions *session.Manager) *Server {
	return &Server{svc: svc, sessions: sessions}
}

type createRequest struct {
	URL string `json:"url"`
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCreateAudit(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req createRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be JSON with a url field.", "")
		return
	}

	a, err := s.svc.Create(r.Context(), req.URL)
	if err != nil {
		s.mapCreateError(w, err)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusCreated, a)
}

func (s *Server) mapCreateError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, validate.ErrEmptyURL),
		errors.Is(err, validate.ErrInvalidURL),
		errors.Is(err, validate.ErrBadScheme),
		errors.Is(err, validate.ErrMissingHost):
		writeError(w, http.StatusBadRequest, "invalid_url", err.Error(), "url")
	default:
		writeError(w, http.StatusInternalServerError, "internal", "Something went wrong.", "")
	}
}

func (s *Server) handleListAudits(w http.ResponseWriter, r *http.Request) {
	if s.sessions != nil && s.sessions.Enabled() && !s.sessions.Authenticated(r, time.Now()) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Sign in to view audit history.", "")
		return
	}
	limit := 20
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			writeError(w, http.StatusBadRequest, "invalid_json", "limit must be an integer between 1 and 100.", "limit")
			return
		}
		limit = n
	}
	list, err := s.svc.List(limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "Something went wrong.", "")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{"items": list})
}

func (s *Server) handleGetAudit(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/audits/")
	id = strings.Trim(id, "/")
	if id == "" {
		writeError(w, http.StatusNotFound, "not_found", "Audit not found.", "")
		return
	}
	a, ok, err := s.svc.Get(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "Something went wrong.", "")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "Audit not found.", "")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, a)
}

type loginRequest struct {
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if s.sessions == nil || !s.sessions.Enabled() {
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "auth": false})
		return
	}
	defer r.Body.Close()
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must be JSON with a password field.", "")
		return
	}
	if !s.sessions.CheckPassword(req.Password) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Incorrect password.", "password")
		return
	}
	s.sessions.SetCookie(w, time.Now())
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if s.sessions != nil {
		s.sessions.ClearCookie(w)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSession(w http.ResponseWriter, r *http.Request) {
	required := s.sessions != nil && s.sessions.Enabled()
	ok := !required || s.sessions.Authenticated(r, time.Now())
	writeJSON(w, http.StatusOK, map[string]any{"authRequired": required, "authenticated": ok})
}
