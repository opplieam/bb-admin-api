package store

import (
	"database/sql"
	"reflect"
	"slices"
	"strings"

	. "github.com/go-jet/jet/v2/postgres"
	. "github.com/opplieam/bb-admin-api/.gen/buy-better-admin/public/table"
	"github.com/opplieam/bb-admin-api/internal/utils"
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

type UnmatchedCategoryDBResult struct {
	TotalRecord int
	ID          int32   `alias:"match_category.id" sql:"primary_key"`
	L1          string  `alias:"match_category.l1"`
	L2          *string `alias:"match_category.l2"`
	L3          *string `alias:"match_category.l3"`
	L4          *string `alias:"match_category.l4"`
	L5          *string `alias:"match_category.l5"`
	L6          *string `alias:"match_category.l6"`
	L7          *string `alias:"match_category.l7"`
	L8          *string `alias:"match_category.l8"`
}

type UnmatchedCategoryResult struct {
	ID            int32  `json:"id"`
	Path          string `json:"path"`
	CategoryLevel int    `json:"category_level"`
}

func (s *CategoryStore) GetUnmatchedCategory(filter utils.Filter) ([]UnmatchedCategoryResult, utils.MetaData, error) {

	stmt := SELECT(
		COUNT(MatchCategory.ID).OVER().AS("UnmatchedCategoryDBResult.TotalRecord"),
		MatchCategory.AllColumns.Except(MatchCategory.MatchID),
	).
		FROM(MatchCategory).
		WHERE(MatchCategory.MatchID.IS_NULL()).
		LIMIT(filter.Limit()).OFFSET(filter.Offset())

	var dest []UnmatchedCategoryDBResult
	if err := stmt.Query(s.DB, &dest); err != nil {
		return nil, utils.MetaData{}, DBTransformError(err)
	}

	var resultList []UnmatchedCategoryResult
	var totalRecord int

	excludeField := []string{"ID", "TotalRecord"}
	for _, v := range dest {
		r := reflect.ValueOf(v)
		var categoryPathList []string
		var categoryLevel int
		var result UnmatchedCategoryResult
		for i := range r.NumField() {
			filedName := r.Type().Field(i).Name
			if slices.Contains(excludeField, filedName) {
				continue
			}
			value := r.Field(i)
			if r.Field(i).Kind() == reflect.Ptr {
				if value.IsNil() {
					continue
				}
				value = value.Elem()
			}
			categoryPathList = append(categoryPathList, value.String())
			categoryLevel++
		}
		result.Path = strings.Join(categoryPathList, " > ")
		result.ID = v.ID
		result.CategoryLevel = categoryLevel
		totalRecord = v.TotalRecord
		resultList = append(resultList, result)
	}

	metaData := utils.GetMetaData(totalRecord, filter.Page, filter.PageSize)

	return resultList, metaData, nil
}
