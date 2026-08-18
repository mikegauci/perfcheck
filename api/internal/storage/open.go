package storage

import "github.com/mikegauci/perfcheck/api/internal/audit"

// Open returns an in-memory repository when path is empty, otherwise SQLite.
func Open(path string) (audit.Repository, error) {
	if path == "" {
		return NewMemory(100), nil
	}
	return NewSQLite(path)
}
