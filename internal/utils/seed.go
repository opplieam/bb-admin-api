package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

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

func SeedCategory(db *sql.DB) error {
	byteData, err := os.ReadFile("data/category.json")
	if err != nil {
		return err
	}
	var data map[string][]string
	err = json.Unmarshal(byteData, &data)
	if err != nil {
		return err
	}

	for k, v := range data {
		parentStmt := Category.INSERT(Category.Name).
			MODEL(model.Category{Name: k}).
			RETURNING(Category.ID)
		parentDest := model.Category{}
		err := parentStmt.Query(db, &parentDest)
		if err != nil {
			return err
		}

		var childModelBulk []model.Category
		for _, vv := range v {
			childModelBulk = append(childModelBulk, model.Category{
				ParentID: parentDest.ID, Name: vv,
			})
		}
		childStmt := Category.INSERT(Category.Name, Category.ParentID).
			MODELS(childModelBulk)
		_, err = childStmt.Exec(db)
		if err != nil {
			return err
		}
	}
	fmt.Println("Seeded Category")
	return nil
}
