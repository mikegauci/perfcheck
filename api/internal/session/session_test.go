package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSignAndValid(t *testing.T) {
	t.Parallel()
	m := New("secret-secret-secret-secret", "pass", false)
	now := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	token := m.Sign(now)
	if !m.Valid(token, now) {
		t.Fatal("fresh token should be valid")
	}
	if m.Valid(token, now.Add(13*time.Hour)) {
		t.Fatal("expired token should be invalid")
	}
	if m.Valid("nope", now) {
		t.Fatal("garbage should be invalid")
	}
}

func TestPasswordCompare(t *testing.T) {
	t.Parallel()
	m := New("secret", "s3cret", false)
	if !m.CheckPassword("s3cret") || m.CheckPassword("nope") {
		t.Fatal("password compare failed")
	}
	open := New("secret", "", false)
	if open.Enabled() || !open.CheckPassword("anything") {
		t.Fatal("empty password should disable auth")
	}
}

func TestCookieRoundTrip(t *testing.T) {
	t.Parallel()
	m := New("secret-secret-secret-secret", "pass", false)
	now := time.Now().UTC()
	rr := httptest.NewRecorder()
	m.SetCookie(rr, now)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range rr.Result().Cookies() {
		req.AddCookie(c)
	}
	if !m.Authenticated(req, now) {
		t.Fatal("expected authenticated request")
	}

	clear := httptest.NewRecorder()
	m.ClearCookie(clear)
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range clear.Result().Cookies() {
		req2.AddCookie(c)
	}
	if m.Authenticated(req2, now) {
		t.Fatal("cleared cookie should not authenticate")
	}
}
