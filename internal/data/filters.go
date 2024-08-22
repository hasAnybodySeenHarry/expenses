package data

import (
	"math"

	"harry2an.com/expenses/internal/validator"
)

type Filters struct {
	Page     int
	PageSize int
}

func ValidateFilters(v *validator.Validator, f *Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
}

func (f *Filters) limit() int {
	return f.PageSize
}

func (f *Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func getMetadata(page, pageSize, total int) Metadata {
	if total == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		TotalRecords: total,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(total) / float64(pageSize))),
	}
}
