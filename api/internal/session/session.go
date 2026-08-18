package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	CookieName = "perfcheck_session"
	ttl        = 12 * time.Hour
)

// Manager issues and verifies HMAC-signed session cookies.
type Manager struct {
	secret   []byte
	password string
	secure   bool
}

func New(secret, password string, secure bool) *Manager {
	if secret == "" {
		var b [32]byte
		_, _ = rand.Read(b[:])
		secret = hex.EncodeToString(b[:])
	}
	return &Manager{secret: []byte(secret), password: password, secure: secure}
}

func (m *Manager) Enabled() bool {
	return m != nil && m.password != ""
}

func (m *Manager) CheckPassword(guess string) bool {
	if !m.Enabled() {
		return true
	}
	return subtle.ConstantTimeCompare([]byte(guess), []byte(m.password)) == 1
}

func (m *Manager) Sign(now time.Time) string {
	exp := now.Add(ttl).Unix()
	msg := strconv.FormatInt(exp, 10)
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(msg))
	return msg + "|" + hex.EncodeToString(mac.Sum(nil))
}

func (m *Manager) Valid(value string, now time.Time) bool {
	parts := strings.Split(value, "|")
	if len(parts) != 2 {
		return false
	}
	exp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false
	}
	if now.Unix() > exp {
		return false
	}
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(parts[0]))
	want := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(parts[1]), []byte(want))
}

func (m *Manager) SetCookie(w http.ResponseWriter, now time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    base64.RawURLEncoding.EncodeToString([]byte(m.Sign(now))),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
		Expires:  now.Add(ttl),
		MaxAge:   int(ttl.Seconds()),
	})
}

func (m *Manager) ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func (m *Manager) Authenticated(r *http.Request, now time.Time) bool {
	if !m.Enabled() {
		return true
	}
	c, err := r.Cookie(CookieName)
	if err != nil {
		return false
	}
	raw, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil {
		return false
	}
	return m.Valid(string(raw), now)
}
