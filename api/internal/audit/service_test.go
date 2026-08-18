package audit

import (
	"context"
	"testing"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/validate"
)

type fakeRepo struct {
	items []Audit
}

func (f *fakeRepo) Save(a Audit) error {
	f.items = append([]Audit{a}, f.items...)
	return nil
}

func (f *fakeRepo) List(limit int) ([]Audit, error) {
	if limit > len(f.items) {
		limit = len(f.items)
	}
	return f.items[:limit], nil
}

func (f *fakeRepo) Get(id string) (Audit, bool, error) {
	for _, a := range f.items {
		if a.ID == id {
			return a, true, nil
		}
	}
	return Audit{}, false, nil
}

var _ Repository = (*fakeRepo)(nil)

func TestServiceCreate(t *testing.T) {
	t.Parallel()
	repo := &fakeRepo{}
	svc := NewService(repo, MockScorer{})
	svc.now = func() time.Time { return time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC) }
	svc.idFn = func() (string, error) { return "fixedid01", nil }

	a, err := svc.Create(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != "fixedid01" || a.URL != "https://example.com" {
		t.Fatalf("unexpected audit: %#v", a)
	}
	if a.Engine != EngineMock {
		t.Fatalf("engine=%q want %q", a.Engine, EngineMock)
	}
	if len(repo.items) != 1 {
		t.Fatalf("expected save, got %d", len(repo.items))
	}
}

func TestServiceCreateInvalid(t *testing.T) {
	t.Parallel()
	svc := NewService(&fakeRepo{}, MockScorer{})
	_, err := svc.Create(context.Background(), "not-a-url")
	if err != validate.ErrInvalidURL {
		t.Fatalf("want ErrInvalidURL, got %v", err)
	}
}
