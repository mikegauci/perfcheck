package audit

// Repository abstracts audit persistence. In-memory for this project;
// a SQLite or Postgres implementation can drop in without changing the service.
type Repository interface {
	Save(a Audit) error
	List(limit int) ([]Audit, error)
	Get(id string) (Audit, bool, error)
}
