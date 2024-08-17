package utils

import (
	"math"
)

type Filter struct {
	Page     int `form:"page,default=1" binding:"min=1,max=1000"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=50"`
}

type MetaData struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func GetMetaData(totalRec, page, pageSize int) MetaData {
	if totalRec == 0 {
		return MetaData{}
	}
	return MetaData{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: totalRec,
		LastPage:     int(math.Ceil(float64(totalRec) / float64(pageSize))),
	}
}

func (f *Filter) Offset() int64 {
	return int64((f.Page - 1) * f.PageSize)
}

func (f *Filter) Limit() int64 {
	return int64(f.PageSize)
}
