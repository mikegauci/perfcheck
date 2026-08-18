package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixtureHandler(t *testing.T, name string, status int, delay time.Duration) http.Handler {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if delay > 0 {
			time.Sleep(delay)
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if name == "clean.html" {
			w.Header().Set("Cache-Control", "public, max-age=60")
		}
		w.WriteHeader(status)
		_, _ = w.Write(body)
	})
}

func TestFetchScorerCleanPage(t *testing.T) {
	t.Parallel()
	srv := httptest.NewTLSServer(fixtureHandler(t, "clean.html", http.StatusOK, 0))
	t.Cleanup(srv.Close)

	scorer := FetchScorer{Client: srv.Client()}
	res, err := scorer.Score(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.Engine != EngineFetch {
		t.Fatalf("engine=%q", res.Engine)
	}
	if res.Signals == nil {
		t.Fatal("expected signals")
	}
	if !res.Signals.HasTitle || !res.Signals.HasMetaDescription || res.Signals.H1Count != 1 {
		t.Fatalf("signals=%#v", res.Signals)
	}
	if res.Signals.ImagesMissingAlt != 0 || !res.Signals.HasLang {
		t.Fatalf("a11y signals=%#v", res.Signals)
	}
	if res.Scores.SEO < 80 {
		t.Fatalf("clean page seo too low: %d", res.Scores.SEO)
	}
}

func TestFetchScorerPoorPage(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(fixtureHandler(t, "poor.html", http.StatusOK, 0))
	t.Cleanup(srv.Close)

	scorer := FetchScorer{Client: srv.Client()}
	res, err := scorer.Score(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.Signals.HasLang {
		t.Fatal("poor page should lack lang")
	}
	if res.Signals.ImagesMissingAlt < 2 {
		t.Fatalf("missing alt=%d", res.Signals.ImagesMissingAlt)
	}
	if res.Scores.Accessibility >= res.Scores.SEO && res.Scores.Accessibility > 80 {
		t.Fatalf("poor page a11y unexpectedly high: %#v", res.Scores)
	}
	if res.Scores.Performance > 95 {
		t.Fatalf("poor page perf unexpectedly high: %d", res.Scores.Performance)
	}
}

func TestFetchScorerNotFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(fixtureHandler(t, "clean.html", http.StatusNotFound, 0))
	t.Cleanup(srv.Close)

	scorer := FetchScorer{Client: srv.Client()}
	res, err := scorer.Score(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.Signals.StatusCode != 404 {
		t.Fatalf("status=%d", res.Signals.StatusCode)
	}
	if res.Scores.Overall >= 80 {
		t.Fatalf("404 should not score well: %#v", res.Scores)
	}
}

func TestFetchScorerSlowPage(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(fixtureHandler(t, "clean.html", http.StatusOK, 450*time.Millisecond))
	t.Cleanup(srv.Close)

	scorer := FetchScorer{Client: srv.Client()}
	res, err := scorer.Score(context.Background(), srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	if res.Signals.TTFBMs < 400 {
		t.Fatalf("ttfb=%d, expected delay", res.Signals.TTFBMs)
	}
}

func TestExtractSignalsFromCleanHTML(t *testing.T) {
	t.Parallel()
	body, err := os.ReadFile(filepath.Join("testdata", "clean.html"))
	if err != nil {
		t.Fatal(err)
	}
	sig := extractSignals(strings.NewReader(string(body)), nil)
	if !sig.HasJSONLD || !sig.HasCanonical || !sig.HasOpenGraph || !sig.HasViewport {
		t.Fatalf("clean extract missing flags: %#v", sig)
	}
	if sig.InputsWithoutLabel != 0 {
		t.Fatalf("clean form should be labelled, got %d", sig.InputsWithoutLabel)
	}
}

func TestScoreFromSignalsDeductions(t *testing.T) {
	t.Parallel()
	perfect := scoreFromSignals(Signals{
		HTTPS: true, HasTitle: true, TitleLength: 30, HasMetaDescription: true,
		H1Count: 1, HasCanonical: true, HasOpenGraph: true, HasJSONLD: true,
		HasLang: true, HasViewport: true, CacheControl: "max-age=60",
	})
	if perfect.Overall != 100 {
		t.Fatalf("perfect=%#v", perfect)
	}
	httpOnly := scoreFromSignals(Signals{
		HTTPS: false, HasTitle: true, TitleLength: 30, HasMetaDescription: true,
		H1Count: 1, HasCanonical: true, HasOpenGraph: true, HasJSONLD: true,
		HasLang: true, HasViewport: true, CacheControl: "max-age=60",
	})
	if httpOnly.Performance >= perfect.Performance {
		t.Fatalf("http should cost performance")
	}
}
