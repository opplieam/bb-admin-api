package utils

import (
	"database/sql"
	"fmt"

	"github.com/opplieam/bb-admin-api/.gen/buy-better-admin/public/model"
	. "github.com/opplieam/bb-admin-api/.gen/buy-better-admin/public/table"
)

func SeedUsers(db *sql.DB) error {
	activeAdmin := model.Users{
		Username: "admin",
		Password: "$2a$10$gzsWRC6/yH2EdNHyXCKnRuO.rEIjMRF/z4GV/a7.hv7alfWdGjZya",
		Active:   true,
	}
	inActiveAdmin := model.Users{
		Username: "pon",
		Password: "$2a$10$3FdE8ZcfpSxSLjla04qvCOY48I718FMLgnJyHLimX1sMvQUcv8aU.",
		Active:   false,
	}
	insertStmt := Users.INSERT(Users.Username, Users.Password, Users.Active).
		MODELS([]model.Users{activeAdmin, inActiveAdmin})
	_, err := insertStmt.Exec(db)
	if err != nil {
		return err
	}
	fmt.Println("Seeded Users")
	return nil
}
