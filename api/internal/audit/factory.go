package audit

import (
	"net/http"
	"time"
)

// ScorerOptions configures the production fetch scorer.
type ScorerOptions struct {
	CacheTTL time.Duration
	HTTP     *http.Client
}

// NewConfiguredScorer builds the production scorer: a live HTML fetch behind a TTL cache.
func NewConfiguredScorer(opts ScorerOptions) Scorer {
	return NewTTLCache(FetchScorer{Client: opts.HTTP}, opts.CacheTTL)
}
