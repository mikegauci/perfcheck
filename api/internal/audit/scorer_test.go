package audit

import (
	"context"
	"testing"
)

func TestMockScorerDeterministic(t *testing.T) {
	t.Parallel()
	url := "https://example.com/pricing"
	scorer := MockScorer{}
	a, err := scorer.Score(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	b, err := scorer.Score(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	if a.Scores != b.Scores || a.Engine != EngineMock {
		t.Fatalf("MockScorer not deterministic: %#v vs %#v", a, b)
	}
}

func TestMockScorerBands(t *testing.T) {
	t.Parallel()
	scorer := MockScorer{}
	urls := []string{
		"https://example.com",
		"https://example.org/a",
		"http://localhost:3000",
		"https://shop.example.com/products?id=1",
	}
	for _, u := range urls {
		res, err := scorer.Score(context.Background(), u)
		if err != nil {
			t.Fatal(err)
		}
		s := res.Scores
		if s.Performance < 45 || s.Performance > 98 {
			t.Fatalf("%s performance=%d out of band", u, s.Performance)
		}
		if s.SEO < 55 || s.SEO > 99 {
			t.Fatalf("%s seo=%d out of band", u, s.SEO)
		}
		if s.Accessibility < 50 || s.Accessibility > 97 {
			t.Fatalf("%s a11y=%d out of band", u, s.Accessibility)
		}
		if s.Overall < 1 || s.Overall > 100 {
			t.Fatalf("%s overall=%d out of range", u, s.Overall)
		}
	}
}

func TestMockScorerDiffersByURL(t *testing.T) {
	t.Parallel()
	scorer := MockScorer{}
	a, _ := scorer.Score(context.Background(), "https://a.example.com")
	b, _ := scorer.Score(context.Background(), "https://b.example.com")
	if a.Scores == b.Scores {
		t.Fatalf("expected different scores for different URLs")
	}
}
