package audit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	fetchTimeout     = 10 * time.Second
	fetchMaxRedirect = 5
	fetchUserAgent   = "PerfCheck/1.0 (+https://github.com/mikegauci/perfcheck)"
)

// FetchScorer scores a URL from one live HTTP GET and HTML inspection.
type FetchScorer struct {
	Client *http.Client
}

func (s FetchScorer) client() *http.Client {
	if s.Client != nil {
		return s.Client
	}
	return &http.Client{
		Timeout: fetchTimeout,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= fetchMaxRedirect {
				return fmt.Errorf("stopped after %d redirects", fetchMaxRedirect)
			}
			return nil
		},
	}
}

func (s FetchScorer) Score(ctx context.Context, pageURL string) (Result, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("User-Agent", fetchUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.8")

	start := time.Now()
	res, err := s.client().Do(req)
	if err != nil {
		return Result{}, err
	}
	defer res.Body.Close()

	ttfb := time.Since(start).Milliseconds()
	limited := io.LimitReader(res.Body, maxBodyBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		return Result{}, err
	}

	base := &Signals{
		StatusCode:      res.StatusCode,
		TTFBMs:          ttfb,
		Bytes:           len(body),
		ContentEncoding: res.Header.Get("Content-Encoding"),
		CacheControl:    res.Header.Get("Cache-Control"),
		HTTPS:           strings.EqualFold(res.Request.URL.Scheme, "https"),
	}
	sig := extractSignals(strings.NewReader(string(body)), base)
	scores := scoreFromSignals(sig)
	return Result{Scores: scores, Engine: EngineFetch, Signals: &sig}, nil
}

func scoreFromSignals(s Signals) Scores {
	perf := 100
	seo := 100
	a11y := 100

	if s.StatusCode >= 400 {
		perf -= 40
		seo -= 30
		a11y -= 20
	}
	if !s.HTTPS {
		perf -= 15
	}
	switch {
	case s.TTFBMs > 800:
		perf -= 20
	case s.TTFBMs > 400:
		perf -= 10
	}
	switch {
	case s.Bytes > 1_500_000:
		perf -= 20
	case s.Bytes > 500_000:
		perf -= 10
	}
	if s.ScriptCount > 15 {
		perf -= 10
	} else if s.ScriptCount > 8 {
		perf -= 5
	}
	if s.StylesheetCount > 8 {
		perf -= 5
	}
	if s.InlineStyleBytes > 20_000 {
		perf -= 10
	}
	if s.CacheControl == "" {
		perf -= 5
	}

	if !s.HasTitle {
		seo -= 25
	} else if s.TitleLength < 15 || s.TitleLength > 70 {
		seo -= 10
	}
	if !s.HasMetaDescription {
		seo -= 20
	}
	if s.H1Count != 1 {
		seo -= 15
	}
	if !s.HasCanonical {
		seo -= 5
	}
	if !s.HasOpenGraph {
		seo -= 5
	}
	if !s.HasJSONLD {
		seo -= 5
	}

	if !s.HasLang {
		a11y -= 20
	}
	if s.ImagesMissingAlt > 0 {
		a11y -= min(30, s.ImagesMissingAlt*5)
	}
	if !s.HasViewport {
		a11y -= 10
	}
	if s.InputsWithoutLabel > 0 {
		a11y -= min(25, s.InputsWithoutLabel*8)
	}

	perf = clampScore(perf)
	seo = clampScore(seo)
	a11y = clampScore(a11y)
	return Scores{
		Performance:   perf,
		SEO:           seo,
		Accessibility: a11y,
		Overall:       (perf*40 + seo*30 + a11y*30) / 100,
	}
}

func clampScore(n int) int {
	if n < 0 {
		return 0
	}
	if n > 100 {
		return 100
	}
	return n
}
