package audit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/validate"
)

// Service orchestrates validation, scoring, recommendations and persistence.
type Service struct {
	repo Repository
	now  func() time.Time
	idFn func() (string, error)
}

// NewService wires a Service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
		idFn: newID,
	}
}

// Create validates the URL, scores it, persists the audit and returns it.
func (s *Service) Create(rawURL string) (Audit, error) {
	normalised, err := validate.NormalizeURL(rawURL)
	if err != nil {
		return Audit{}, err
	}

	id, err := s.idFn()
	if err != nil {
		return Audit{}, fmt.Errorf("generate id: %w", err)
	}

	scores := ScoreURL(normalised)
	a := Audit{
		ID:              id,
		URL:             normalised,
		CreatedAt:       s.now(),
		Scores:          scores,
		Recommendations: Recommend(scores),
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
