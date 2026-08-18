package audit

import "testing"

func TestScoreURLDeterministic(t *testing.T) {
	t.Parallel()
	url := "https://example.com/pricing"
	a := ScoreURL(url)
	b := ScoreURL(url)
	if a != b {
		t.Fatalf("ScoreURL not deterministic: %#v vs %#v", a, b)
	}
}

func TestScoreURLBands(t *testing.T) {
	t.Parallel()
	urls := []string{
		"https://example.com",
		"https://example.org/a",
		"http://localhost:3000",
		"https://shop.example.com/products?id=1",
	}
	for _, u := range urls {
		s := ScoreURL(u)
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

func TestScoreURLDiffersByURL(t *testing.T) {
	t.Parallel()
	a := ScoreURL("https://a.example.com")
	b := ScoreURL("https://b.example.com")
	if a == b {
		t.Fatalf("expected different scores for different URLs")
	}
}
