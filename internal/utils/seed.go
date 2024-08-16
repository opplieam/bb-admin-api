package utils

import (
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/go-faker/faker/v4"
	. "github.com/go-jet/jet/v2/postgres"
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
	maxChildLevel = 4
)

func SeedCategory(db *sql.DB) error {
	fmt.Println("Seeded Category")
	return insertCategoryRecursive(db, nil, 1)
}

func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}

func insertCategoryRecursive(db *sql.DB, parentID *int32, level int) error {
	var maxCategory int
	if parentID == nil {
		maxCategory = 5
	} else {
		maxCategory = randRange(1, 3)
	}
	for i := 1; i <= maxCategory; i++ {
		catModel := model.Category{
			Name:     faker.Word(),
			ParentID: parentID,
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

		if err := insertCategoryRecursive(db, &childDest.ID, level+1); err != nil {
			return err
		}
	}
	return nil
}

func getValueOrNil(index int, value []string, totalLevel int) *string {
	if index > totalLevel {
		return nil
	}
	return &value[index-1]
}

func SeedMatchCategory(db *sql.DB) error {
	var totalRecords = 102
	column := map[int]ColumnString{
		0: MatchCategory.L1,
		1: MatchCategory.L2,
		2: MatchCategory.L3,
		3: MatchCategory.L4,
		4: MatchCategory.L5,
		5: MatchCategory.L6,
		6: MatchCategory.L7,
		7: MatchCategory.L8,
	}

	for range totalRecords {
		totalLevel := randRange(2, 8)
		var selectColumn ColumnList
		var value []string
		for i := range totalLevel {
			selectColumn = append(selectColumn, column[i])
			value = append(value, faker.Word())
		}
		stmt := MatchCategory.INSERT(selectColumn).MODEL(model.MatchCategory{
			L1: faker.Word(),
			L2: getValueOrNil(1, value, totalLevel),
			L3: getValueOrNil(2, value, totalLevel),
			L4: getValueOrNil(3, value, totalLevel),
			L5: getValueOrNil(4, value, totalLevel),
			L6: getValueOrNil(5, value, totalLevel),
			L7: getValueOrNil(6, value, totalLevel),
			L8: getValueOrNil(7, value, totalLevel),
		})

		if _, err := stmt.Exec(db); err != nil {
			return err
		}
	}
	fmt.Println("Seeded Match Category")
	return nil
}
