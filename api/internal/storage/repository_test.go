package storage

import (
	"testing"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
)

func sampleAudit(id string) audit.Audit {
	return audit.Audit{
		ID:        id,
		URL:       "https://example.com/" + id,
		CreatedAt: time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC).Add(time.Duration(id[len(id)-1]) * time.Second),
		Engine:    audit.EngineMock,
		Scores:    audit.Scores{Overall: 70, Performance: 70, SEO: 70, Accessibility: 70},
		Recommendations: []audit.Recommendation{
			{ID: "img-format", Category: "performance", Severity: "high", Title: "Images", Detail: "Use AVIF"},
		},
	}
}

func TestRepositoryConformance(t *testing.T) {
	t.Run("memory", func(t *testing.T) {
		assertRepository(t, NewMemory(100))
	})
	t.Run("sqlite", func(t *testing.T) {
		repo, err := NewSQLite(t.TempDir() + "/perfcheck.db")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = repo.Close() })
		assertRepository(t, repo)
	})
}

func assertRepository(t *testing.T, repo audit.Repository) {
	t.Helper()
	a := sampleAudit("one")
	if err := repo.Save(a); err != nil {
		t.Fatal(err)
	}
	got, ok, err := repo.Get(a.ID)
	if err != nil || !ok {
		t.Fatalf("get: ok=%v err=%v", ok, err)
	}
	if got.URL != a.URL || got.Engine != a.Engine || got.Scores.Overall != 70 {
		t.Fatalf("round trip mismatch: %#v", got)
	}
	if len(got.Recommendations) != 1 {
		t.Fatalf("recommendations=%v", got.Recommendations)
	}

	b := sampleAudit("two")
	if err := repo.Save(b); err != nil {
		t.Fatal(err)
	}
	list, err := repo.List(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) < 2 {
		t.Fatalf("list len=%d", len(list))
	}
	if _, ok, err := repo.Get("missing"); err != nil || ok {
		t.Fatalf("missing should be not found, ok=%v err=%v", ok, err)
	}
}

func TestOpenMemoryDefault(t *testing.T) {
	repo, err := Open("")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := repo.(*Memory); !ok {
		t.Fatalf("want Memory, got %T", repo)
	}
}
