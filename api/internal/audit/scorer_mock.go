package audit

import (
	"context"
	"hash/fnv"
	"math/rand"
)

// MockScorer derives deterministic scores from a URL hash.
// This is intentionally synthetic — there is no network call.
type MockScorer struct{}

func (MockScorer) Score(_ context.Context, normalisedURL string) (Result, error) {
	return Result{
		Scores: scoreURL(normalisedURL),
		Engine: EngineMock,
	}, nil
}

func scoreURL(normalisedURL string) Scores {
	h := fnv.New64a()
	_, _ = h.Write([]byte(normalisedURL))
	seed := h.Sum64()
	r := rand.New(rand.NewSource(int64(seed))) //nolint:gosec // deterministic mock, not crypto

	performance := band(r, 45, 98)
	seo := band(r, 55, 99)
	accessibility := band(r, 50, 97)
	overall := (performance*40 + seo*30 + accessibility*30) / 100

	return Scores{
		Overall:       overall,
		Performance:   performance,
		SEO:           seo,
		Accessibility: accessibility,
	}
}

func band(r *rand.Rand, min, max int) int {
	if max <= min {
		return min
	}
	return min + r.Intn(max-min+1)
}
