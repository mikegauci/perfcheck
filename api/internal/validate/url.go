package validate

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrEmptyURL    = errors.New("url is required")
	ErrInvalidURL  = errors.New("enter a full URL including https://")
	ErrBadScheme   = errors.New("only http and https URLs are supported")
	ErrMissingHost = errors.New("URL must include a hostname")
)

// NormalizeURL parses and normalises a user-supplied URL.
// It requires an absolute URL with http or https scheme and a host.
func NormalizeURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrEmptyURL
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", ErrInvalidURL
	}
	if u.Scheme == "" || u.Host == "" {
		return "", ErrInvalidURL
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return "", ErrBadScheme
	}
	if strings.TrimSpace(u.Hostname()) == "" {
		return "", ErrMissingHost
	}

	u.Scheme = scheme
	u.Fragment = ""
	// Drop default ports for stable hashing.
	if (scheme == "http" && u.Port() == "80") || (scheme == "https" && u.Port() == "443") {
		u.Host = u.Hostname()
	}

	return u.String(), nil
}
