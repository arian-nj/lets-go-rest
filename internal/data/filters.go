package data

import "github.com/arian-nj/site/back/internal/validator"

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilter(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be bigger than 0")
	v.Check(f.Page <= 10_000_000, "page", "must be smaller than 10 milion")
	v.Check(f.PageSize > 0, "page_size", "must be bigger than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be smaller than 100")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
