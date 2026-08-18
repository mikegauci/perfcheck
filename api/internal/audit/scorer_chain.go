package audit

import (
	"context"
	"errors"
	"log"
)

// ChainScorer tries each scorer in order and returns the first success.
type ChainScorer struct {
	Scorers []Scorer
}

func (c ChainScorer) Score(ctx context.Context, pageURL string) (Result, error) {
	if len(c.Scorers) == 0 {
		return Result{}, errors.New("no scorers configured")
	}
	var last error
	for i, scorer := range c.Scorers {
		res, err := scorer.Score(ctx, pageURL)
		if err == nil {
			if i > 0 {
				log.Printf("scorer fallback succeeded with engine %s after: %v", res.Engine, last)
			}
			return res, nil
		}
		last = err
		log.Printf("scorer %d failed: %v", i, err)
	}
	return Result{}, last
}
