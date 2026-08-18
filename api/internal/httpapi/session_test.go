package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/session"
	"github.com/mikegauci/perfcheck/api/internal/storage"
)

func TestListRequiresSessionWhenConfigured(t *testing.T) {
	t.Parallel()
	repo := storage.NewMemory(10)
	svc := audit.NewService(repo, stubScorer{})
	mgr := session.New("unit-test-secret-unit-test-secret", "letmein", false)
	router := NewRouter(Options{Service: svc, Sessions: mgr, CORSOrigin: "http://localhost:1313"})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/audits", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rr.Code)
	}

	loginBody := bytes.NewBufferString(`{"password":"letmein"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/v1/session", loginBody)
	loginRR := httptest.NewRecorder()
	router.ServeHTTP(loginRR, loginReq)
	if loginRR.Code != http.StatusOK {
		t.Fatalf("login status=%d body=%s", loginRR.Code, loginRR.Body.String())
	}

	authed := httptest.NewRequest(http.MethodGet, "/api/v1/audits", nil)
	for _, c := range loginRR.Result().Cookies() {
		authed.AddCookie(c)
	}
	authedRR := httptest.NewRecorder()
	router.ServeHTTP(authedRR, authed)
	if authedRR.Code != http.StatusOK {
		t.Fatalf("authed list status=%d", authedRR.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(authedRR.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
}
