package audit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestPSIScorerFixture(t *testing.T) {
	t.Parallel()
	body, err := os.ReadFile(filepath.Join("testdata", "psi.json"))
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("url") == "" {
			t.Error("missing url param")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	scorer := PSIScorer{Client: srv.Client(), Base: srv.URL}
	res, err := scorer.Score(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if res.Engine != EnginePSI {
		t.Fatalf("engine=%s", res.Engine)
	}
	if res.Scores.Performance != 72 || res.Scores.SEO != 91 || res.Scores.Accessibility != 84 {
		t.Fatalf("scores=%#v", res.Scores)
	}
	if res.Signals == nil || res.Signals.ImagesMissingAlt != 1 {
		t.Fatalf("signals=%#v", res.Signals)
	}
}

func TestPSIScorerHTTPError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusTooManyRequests)
	}))
	t.Cleanup(srv.Close)
	scorer := PSIScorer{Client: srv.Client(), Base: srv.URL}
	_, err := scorer.Score(context.Background(), "https://example.com")
	if err == nil {
		t.Fatal("expected error")
	}
}

type seqScorer struct {
	engine string
	err    error
	calls  *int32
}

func (s seqScorer) Score(context.Context, string) (Result, error) {
	if s.calls != nil {
		atomic.AddInt32(s.calls, 1)
	}
	if s.err != nil {
		return Result{}, s.err
	}
	return Result{Engine: s.engine, Scores: Scores{Overall: 50, Performance: 50, SEO: 50, Accessibility: 50}}, nil
}

func TestChainScorerFallback(t *testing.T) {
	t.Parallel()
	var psiCalls, fetchCalls int32
	chain := ChainScorer{Scorers: []Scorer{
		seqScorer{engine: EnginePSI, err: context.DeadlineExceeded, calls: &psiCalls},
		seqScorer{engine: EngineFetch, calls: &fetchCalls},
		MockScorer{},
	}}
	res, err := chain.Score(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if res.Engine != EngineFetch {
		t.Fatalf("engine=%s", res.Engine)
	}
	if psiCalls != 1 || fetchCalls != 1 {
		t.Fatalf("calls psi=%d fetch=%d", psiCalls, fetchCalls)
	}
}

func TestTTLCacheHitAndExpiry(t *testing.T) {
	t.Parallel()
	var calls int32
	inner := seqScorer{engine: EngineFetch, calls: &calls}
	cache := NewTTLCache(inner, 40*time.Millisecond)

	if _, err := cache.Score(context.Background(), "https://a.example"); err != nil {
		t.Fatal(err)
	}
	if _, err := cache.Score(context.Background(), "https://a.example"); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected cache hit, calls=%d", calls)
	}
	time.Sleep(50 * time.Millisecond)
	if _, err := cache.Score(context.Background(), "https://a.example"); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected expiry, calls=%d", calls)
	}
}

func TestNewConfiguredScorerModes(t *testing.T) {
	t.Parallel()
	s, err := NewConfiguredScorer(ScorerOptions{Mode: "mock", CacheTTL: time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	res, err := s.Score(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if res.Engine != EngineMock {
		t.Fatalf("engine=%s", res.Engine)
	}
	if _, err := NewConfiguredScorer(ScorerOptions{Mode: "nope"}); err == nil {
		t.Fatal("expected unknown mode error")
	}
}
