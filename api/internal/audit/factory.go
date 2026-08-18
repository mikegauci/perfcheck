package audit

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

// ScorerOptions selects which engines to run.
type ScorerOptions struct {
	Mode     string // auto | psi | fetch | mock
	PSIKey   string
	CacheTTL time.Duration
	HTTP     *http.Client
}

// NewConfiguredScorer builds the production scorer chain.
func NewConfiguredScorer(opts ScorerOptions) (Scorer, error) {
	mode := opts.Mode
	if mode == "" {
		mode = "auto"
	}
	key := opts.PSIKey
	if key == "" {
		key = os.Getenv("PSI_API_KEY")
	}

	var inner Scorer
	switch mode {
	case "mock":
		inner = MockScorer{}
	case "fetch":
		inner = ChainScorer{Scorers: []Scorer{FetchScorer{Client: opts.HTTP}, MockScorer{}}}
	case "psi":
		inner = ChainScorer{Scorers: []Scorer{
			PSIScorer{Client: opts.HTTP, APIKey: key},
			FetchScorer{Client: opts.HTTP},
			MockScorer{},
		}}
	case "auto":
		inner = ChainScorer{Scorers: []Scorer{
			PSIScorer{Client: opts.HTTP, APIKey: key},
			FetchScorer{Client: opts.HTTP},
			MockScorer{},
		}}
	default:
		return nil, fmt.Errorf("unknown scorer mode %q", mode)
	}

	return NewTTLCache(inner, opts.CacheTTL), nil
}
