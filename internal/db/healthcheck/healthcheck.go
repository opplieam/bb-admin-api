package healthcheck

import (
	"database/sql"
)

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func (s *Store) HealthCheck() (bool, error) {
	tmp := false
	rows, err := s.DB.Query(`SELECT true`)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		err = rows.Scan(&tmp)
		if err != nil {
			return false, err
		}
	}

	return tmp, nil

}
