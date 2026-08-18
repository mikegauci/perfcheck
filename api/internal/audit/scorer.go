package audit

import (
	"hash/fnv"
	"math/rand"
)

// ScoreURL derives deterministic mock scores from a normalised URL.
// This is intentionally synthetic — there is no real Lighthouse call.
func ScoreURL(normalisedURL string) Scores {
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
