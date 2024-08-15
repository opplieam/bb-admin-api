package store

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	. "github.com/opplieam/bb-admin-api/.gen/buy-better-admin/public/table"
)

type CategoryStore struct {
	DB *sql.DB
}

func NewCategoryStore(db *sql.DB) *CategoryStore {
	return &CategoryStore{
		DB: db,
	}
}

type AllCategoryResult struct {
	ID       int32  `sql:"primary_key" alias:"category.id" json:"id"`
	Name     string `alias:"category.name" json:"name"`
	HasChild bool   `alias:"category.has_child" json:"has_child"`
	Path     string `json:"path"`
}

func (s *CategoryStore) GetAllCategory() ([]AllCategoryResult, error) {
	catRecur := CTE("catRecur")
	pathCol := StringColumn("AllCategoryResult.path").From(catRecur)
	stmt := WITH_RECURSIVE(
		catRecur.AS(
			SELECT(
				Category.ID, Category.Name, Category.HasChild, CAST(Category.Name).AS_TEXT().AS(pathCol.Name()),
			).FROM(
				Category,
			).WHERE(
				Category.ParentID.IS_NULL(),
			).UNION(
				SELECT(
					Category.ID, Category.Name, Category.HasChild, CAST(pathCol.CONCAT(String(" > ")).CONCAT(Category.Name)).AS_TEXT(),
				).FROM(
					Category.
						INNER_JOIN(catRecur, Category.ParentID.EQ(Category.ID.From(catRecur))),
				),
			),
		),
	)(
		SELECT(
			catRecur.AllColumns(),
		).FROM(
			catRecur,
		).WHERE(
			Category.HasChild.From(catRecur).IS_FALSE(),
		),
	)
	var dest []AllCategoryResult
	if err := stmt.Query(s.DB, &dest); err != nil {
		return nil, DBTransformError(err)
	}

	return dest, nil
}
