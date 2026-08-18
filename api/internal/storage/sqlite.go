package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mikegauci/perfcheck/api/internal/audit"
	_ "modernc.org/sqlite"
)

// SQLite is an audit.Repository backed by a local database file.
type SQLite struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	if path == "" {
		return nil, fmt.Errorf("sqlite path is required")
	}
	dsn := path
	if !strings.Contains(path, "?") {
		dsn = path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)"
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	s := &SQLite{db: db}
	if err := s.migrate(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) migrate() error {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS audits (
  id          TEXT PRIMARY KEY,
  url         TEXT NOT NULL,
  created_at  TEXT NOT NULL,
  engine      TEXT NOT NULL,
  overall     INTEGER NOT NULL,
  payload     TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_audits_created_at ON audits(created_at DESC);
`)
	return err
}

var _ audit.Repository = (*SQLite)(nil)

func (s *SQLite) Save(a audit.Audit) error {
	payload, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		`INSERT INTO audits (id, url, created_at, engine, overall, payload)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   url=excluded.url,
		   created_at=excluded.created_at,
		   engine=excluded.engine,
		   overall=excluded.overall,
		   payload=excluded.payload`,
		a.ID,
		a.URL,
		a.CreatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"),
		a.Engine,
		a.Scores.Overall,
		string(payload),
	)
	return err
}

func (s *SQLite) List(limit int) ([]audit.Audit, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(
		`SELECT payload FROM audits ORDER BY created_at DESC LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]audit.Audit, 0, limit)
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, err
		}
		var a audit.Audit
		if err := json.Unmarshal([]byte(payload), &a); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *SQLite) Get(id string) (audit.Audit, bool, error) {
	var payload string
	err := s.db.QueryRow(`SELECT payload FROM audits WHERE id = ?`, id).Scan(&payload)
	if err == sql.ErrNoRows {
		return audit.Audit{}, false, nil
	}
	if err != nil {
		return audit.Audit{}, false, err
	}
	var a audit.Audit
	if err := json.Unmarshal([]byte(payload), &a); err != nil {
		return audit.Audit{}, false, err
	}
	return a, true, nil
}
