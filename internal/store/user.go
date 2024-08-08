package store

import (
	"database/sql"
	"fmt"
	"time"

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
		return 0, DBTransformError(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dest.Password), []byte(password)); err != nil {
		return 0, fmt.Errorf("password mismatch, %w", err)
	}

	return dest.ID, nil
}

// CreateUser Create admin user, return nil if success
func (s *UserStore) CreateUser(username, password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	user := model.Users{
		Username: username,
		Password: string(hashPassword),
	}

	stmt := Users.INSERT(Users.Username, Users.Password).MODEL(user)
	_, err = stmt.Exec(s.DB)
	if err != nil {
		return DBTransformError(err)

	}

	return nil
}

type AllUsersResult struct {
	ID        int32     `sql:"primary_key" alias:"users.id" json:"id"`
	CreatedAt time.Time `alias:"users.created_at" json:"created_at"`
	UpdatedAt time.Time `alias:"users.updated_at" json:"updated_at"`
	Username  string    `alias:"users.username" json:"username"`
	Active    bool      `alias:"users.active" json:"active"`
}

// GetAllUsers return list of users
func (s *UserStore) GetAllUsers() ([]AllUsersResult, error) {
	stmt := SELECT(
		Users.ID, Users.Username, Users.CreatedAt, Users.UpdatedAt, Users.Active,
	).FROM(
		Users,
	).ORDER_BY(
		Users.ID,
	)

	var dest []AllUsersResult
	if err := stmt.Query(s.DB, &dest); err != nil {
		return nil, DBTransformError(err)
	}
	return dest, nil
}

// UpdateUserStatus Update user status, return nil if success
func (s *UserStore) UpdateUserStatus(userId int32, active bool) error {
	userModel := model.Users{Active: active, UpdatedAt: time.Now()}
	stmt := Users.
		UPDATE(Users.Active, Users.UpdatedAt).
		MODEL(userModel).
		WHERE(Users.ID.EQ(Int32(userId)))

	_, err := stmt.Exec(s.DB)
	if err != nil {
		return DBTransformError(err)
	}
	return nil
}
