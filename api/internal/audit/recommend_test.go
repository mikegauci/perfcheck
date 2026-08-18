package audit

import "testing"

func TestRecommendEmptyWhenPerfect(t *testing.T) {
	t.Parallel()
	recs := Recommend(Result{Scores: Scores{Overall: 90, Performance: 95, SEO: 95, Accessibility: 95}})
	if len(recs) == 0 {
		t.Fatal("expected at least the overall-pass recommendation")
	}
}

func TestRecommendLowPerformance(t *testing.T) {
	t.Parallel()
	recs := Recommend(Result{Scores: Scores{Overall: 60, Performance: 50, SEO: 95, Accessibility: 95}})
	found := false
	for _, r := range recs {
		if r.ID == "img-format" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected img-format recommendation for low performance")
	}
}

func TestRecommendUsesSignals(t *testing.T) {
	t.Parallel()
	recs := Recommend(Result{
		Scores:  Scores{Overall: 95, Performance: 95, SEO: 95, Accessibility: 95},
		Signals: &Signals{ImagesMissingAlt: 3, HasLang: true, HasMetaDescription: true, H1Count: 1},
	})
	found := false
	for _, r := range recs {
		if r.ID == "missing-alt" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected missing-alt when images lack alt text")
	}
}
