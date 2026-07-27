
package pagination

import (
	"math"
)

type Pagination struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
	NextPage   string      `json:"next_page,omitempty"`
	PrevPage   string      `json:"prev_page,omitempty"`
}



// a method to create pagination that collects page and limit from query params
func CreatePagination(page, limit int, totalRows int64, data interface{}) Pagination {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       data,
		NextPage: GetNextPageURL(nil, page, limit, totalPages),
		PrevPage: GetPrevPageURL(nil, page, limit),
	}
}

 

