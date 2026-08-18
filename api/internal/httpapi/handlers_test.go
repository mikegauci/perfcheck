package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	"github.com/mikegauci/perfcheck/api/internal/storage"
)

func testRouter() http.Handler {
	repo := storage.NewMemory(50)
	svc := audit.NewService(repo)
	return NewRouter(Options{Service: svc, CORSOrigin: "http://localhost:1313"})
}

func TestHealthz(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/healthz", nil)
	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d", rr.Code)
	}
}

func TestCreateAndGetAudit(t *testing.T) {
	t.Parallel()
	router := testRouter()

	body := bytes.NewBufferString(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/audits", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rr.Code, rr.Body.String())
	}

	var created audit.Audit
	if err := json.Unmarshal(rr.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.Scores.Overall == 0 {
		t.Fatalf("unexpected audit: %#v", created)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/audits/"+created.ID, nil)
	getRR := httptest.NewRecorder()
	router.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("get status=%d", getRR.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/audits?limit=10", nil)
	listRR := httptest.NewRecorder()
	router.ServeHTTP(listRR, listReq)
	if listRR.Code != http.StatusOK {
		t.Fatalf("list status=%d", listRR.Code)
	}
}

func TestCreateInvalidURL(t *testing.T) {
	t.Parallel()
	body := bytes.NewBufferString(`{"url":"not-a-url"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/audits", body)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rr.Code)
	}
	var env ErrorBody
	if err := json.Unmarshal(rr.Body.Bytes(), &env); err != nil {
		t.Fatal(err)
	}
	if env.Error.Code != "invalid_url" || env.Error.Field != "url" {
		t.Fatalf("envelope=%#v", env)
	}
}

func TestGetNotFound(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/audits/missing", nil)
	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rr.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/audits", nil)
	req.Header.Set("Origin", "http://localhost:1313")
	rr := httptest.NewRecorder()
	testRouter().ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status=%d", rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "http://localhost:1313" {
		t.Fatalf("missing CORS header")
	}
}
