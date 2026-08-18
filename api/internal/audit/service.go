package audit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/validate"
)

// Service orchestrates validation, scoring, recommendations and persistence.
type Service struct {
	repo   Repository
	scorer Scorer
	now    func() time.Time
	idFn   func() (string, error)
}

// NewService wires a Service with the given repository and scorer.
func NewService(repo Repository, scorer Scorer) *Service {
	if scorer == nil {
		scorer = FetchScorer{}
	}
	return &Service{
		repo:   repo,
		scorer: scorer,
		now:    func() time.Time { return time.Now().UTC() },
		idFn:   newID,
	}
}

// Create validates the URL, scores it, persists the audit and returns it.
func (s *Service) Create(ctx context.Context, rawURL string) (Audit, error) {
	normalised, err := validate.NormalizeURL(rawURL)
	if err != nil {
		return Audit{}, err
	}

	id, err := s.idFn()
	if err != nil {
		return Audit{}, fmt.Errorf("generate id: %w", err)
	}

	res, err := s.scorer.Score(ctx, normalised)
	if err != nil {
		return Audit{}, fmt.Errorf("score url: %w", err)
	}

	a := Audit{
		ID:              id,
		URL:             normalised,
		CreatedAt:       s.now(),
		Engine:          res.Engine,
		Scores:          res.Scores,
		Signals:         res.Signals,
		Recommendations: Recommend(res),
	}

	if err := s.repo.Save(a); err != nil {
		return Audit{}, fmt.Errorf("save audit: %w", err)
	}
	return a, nil
}

// List returns the newest audits up to limit.
func (s *Service) List(limit int) ([]Audit, error) {
	return s.repo.List(limit)
}

// Get returns a single audit by id.
func (s *Service) Get(id string) (Audit, bool, error) {
	return s.repo.Get(id)
}

func newID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
