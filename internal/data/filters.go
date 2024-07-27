package data

import (
	"strings"

	"github.com/arian-nj/site/back/internal/validator"
)

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

func (f *Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f *Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
