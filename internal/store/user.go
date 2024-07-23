package store

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/opplieam/bb-admin-api/.gen/buy-better-admin/public/model"
	. "github.com/opplieam/bb-admin-api/.gen/buy-better-admin/public/table"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	DB *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{DB: db}
}

// IsAuthenticated check username and hashed password
// Parameters:
//
//	username (string): username
//	password (string): plain password
//
// Returns:
//
//	error: return nil if it authenticated
func (s *UserStore) IsAuthenticated(username, password string) error {
	stmt := SELECT(
		Users.Password,
	).FROM(
		Users,
	).WHERE(
		Users.Username.EQ(String(username)).
			AND(Users.Active.IS_TRUE()),
	)
	var dest struct {
		model.Users
	}
	if err := stmt.Query(s.DB, &dest); err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dest.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}
