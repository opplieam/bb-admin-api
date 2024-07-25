package store

import (
	"database/sql"
	"fmt"

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

// FindByCredential return user_id if credential match
func (s *UserStore) FindByCredential(username, password string) (int32, error) {
	stmt := SELECT(
		Users.ID, Users.Password,
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
		return 0, ErrRecordNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dest.Password), []byte(password)); err != nil {
		return 0, fmt.Errorf("password mismatch, %w", err)
	}

	return dest.ID, nil
}
