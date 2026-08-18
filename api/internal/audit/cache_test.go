package audit

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type stubScorer struct {
	engine string
	err    error
	calls  *int32
	scores Scores
}

func (s stubScorer) Score(context.Context, string) (Result, error) {
	if s.calls != nil {
		atomic.AddInt32(s.calls, 1)
	}
	if s.err != nil {
		return Result{}, s.err
	}
	engine := s.engine
	if engine == "" {
		engine = EngineFetch
	}
	scores := s.scores
	if scores == (Scores{}) {
		scores = Scores{Overall: 80, Performance: 80, SEO: 80, Accessibility: 80}
	}
	return Result{Engine: engine, Scores: scores}, nil
}

func TestTTLCacheHitAndExpiry(t *testing.T) {
	t.Parallel()
	var calls int32
	inner := stubScorer{engine: EngineFetch, calls: &calls}
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

func TestNewConfiguredScorer(t *testing.T) {
	t.Parallel()
	s := NewConfiguredScorer(ScorerOptions{CacheTTL: time.Minute})
	if s == nil {
		t.Fatal("expected scorer")
	}
	cache, ok := s.(*TTLCache)
	if !ok {
		t.Fatalf("want *TTLCache, got %T", s)
	}
	if _, ok := cache.inner.(FetchScorer); !ok {
		t.Fatalf("want FetchScorer inner, got %T", cache.inner)
	}
}
