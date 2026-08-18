package storage

import (
	"sync"

	"github.com/mikegauci/perfcheck/api/internal/audit"
)

const defaultMaxHistory = 100

// Memory is a mutex-guarded in-memory audit.Repository with a capped history.
type Memory struct {
	mu      sync.RWMutex
	items   []audit.Audit
	byID    map[string]int
	maxSize int
}

// NewMemory returns an in-memory store capped at maxSize entries (newest kept).
func NewMemory(maxSize int) *Memory {
	if maxSize <= 0 {
		maxSize = defaultMaxHistory
	}
	return &Memory{
		items:   make([]audit.Audit, 0, maxSize),
		byID:    make(map[string]int),
		maxSize: maxSize,
	}
}

var _ audit.Repository = (*Memory)(nil)

func (m *Memory) Save(a audit.Audit) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if idx, ok := m.byID[a.ID]; ok {
		m.items[idx] = a
		return nil
	}

	m.items = append([]audit.Audit{a}, m.items...)
	m.reindexLocked()

	if len(m.items) > m.maxSize {
		m.items = m.items[:m.maxSize]
		m.reindexLocked()
	}
	return nil
}

func (m *Memory) List(limit int) ([]audit.Audit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit <= 0 || limit > len(m.items) {
		limit = len(m.items)
	}
	out := make([]audit.Audit, limit)
	copy(out, m.items[:limit])
	return out, nil
}

func (m *Memory) Get(id string) (audit.Audit, bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	idx, ok := m.byID[id]
	if !ok {
		return audit.Audit{}, false, nil
	}
	return m.items[idx], true, nil
}

func (m *Memory) reindexLocked() {
	m.byID = make(map[string]int, len(m.items))
	for i, a := range m.items {
		m.byID[a.ID] = i
	}
}
