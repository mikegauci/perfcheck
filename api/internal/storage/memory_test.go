package storage

import (
	"sync"
	"testing"
	"time"

	"github.com/mikegauci/perfcheck/api/internal/audit"
)

func sample(id string) audit.Audit {
	return audit.Audit{
		ID:        id,
		URL:       "https://example.com/" + id,
		CreatedAt: time.Now().UTC(),
		Scores:    audit.Scores{Overall: 70, Performance: 70, SEO: 70, Accessibility: 70},
	}
}

func TestMemorySaveListGet(t *testing.T) {
	t.Parallel()
	repo := NewMemory(10)
	a := sample("a1")
	if err := repo.Save(a); err != nil {
		t.Fatal(err)
	}
	got, ok, err := repo.Get("a1")
	if err != nil || !ok || got.ID != "a1" {
		t.Fatalf("Get failed: ok=%v err=%v got=%v", ok, err, got)
	}
	list, err := repo.List(5)
	if err != nil || len(list) != 1 {
		t.Fatalf("List = %v err=%v", list, err)
	}
}

func TestMemoryCap(t *testing.T) {
	t.Parallel()
	repo := NewMemory(2)
	_ = repo.Save(sample("1"))
	_ = repo.Save(sample("2"))
	_ = repo.Save(sample("3"))
	list, err := repo.List(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("len=%d want 2", len(list))
	}
	if list[0].ID != "3" || list[1].ID != "2" {
		t.Fatalf("order wrong: %#v", list)
	}
	if _, ok, _ := repo.Get("1"); ok {
		t.Fatal("expected oldest to be dropped")
	}
}

func TestMemoryConcurrentWrites(t *testing.T) {
	t.Parallel()
	repo := NewMemory(100)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = repo.Save(sample(string(rune('a'+n%26)) + string(rune('0'+n%10))))
		}(i)
	}
	wg.Wait()
	list, err := repo.List(100)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) == 0 {
		t.Fatal("expected some saves")
	}
}
