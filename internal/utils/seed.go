package utils

import (
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
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

var (
	maxNumTopCategory = 5
	maxChildLevel     = 4
)

func SeedCategory(db *sql.DB) error {
	for i := 0; i < maxNumTopCategory; i++ {
		if err := insertCategoryRecursive(db, -1, 1); err != nil {
			return err
		}
	}
	fmt.Println("Seeded Category")
	return nil
}

func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

func insertCategoryRecursive(db *sql.DB, parentID int32, level int) error {
	if parentID == -1 {
		parentStmt := Category.INSERT(Category.Name).
			MODEL(model.Category{Name: faker.Word()}).
			RETURNING(Category.ID)
		var parentDest model.Category
		if err := parentStmt.Query(db, &parentDest); err != nil {
			return err
		}
		return insertCategoryRecursive(db, parentDest.ID, level+1)
	}

	for i := 1; i <= randRange(1, 3); i++ {
		catModel := model.Category{
			Name:     faker.Word(),
			ParentID: &parentID,
			HasChild: level != maxChildLevel,
		}
		stmt := Category.INSERT(Category.Name, Category.ParentID, Category.HasChild).
			MODEL(catModel).
			RETURNING(Category.ID)
		var childDest model.Category
		if err := stmt.Query(db, &childDest); err != nil {
			return err
		}

		if level == maxChildLevel {
			continue
		}

		if err := insertCategoryRecursive(db, childDest.ID, level+1); err != nil {
			return err
		}
	}
	return nil
}
