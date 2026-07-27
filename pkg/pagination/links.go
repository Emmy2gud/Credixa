package pagination

import (
	"net/http"
	"strconv"
)

func GetNextPageURL(r *http.Request, page, limit, totalPages int) string {
	if page >= totalPages {
		return ""
	}

	q := r.URL.Query()
	q.Set("page", strconv.Itoa(page+1))
	q.Set("limit", strconv.Itoa(limit))

	return r.URL.Path + "?" + q.Encode()
}

func GetPrevPageURL(r *http.Request, page, limit int) string {
	if page <= 1 {
		return ""
	}

	q := r.URL.Query()
	q.Set("page", strconv.Itoa(page-1))
	q.Set("limit", strconv.Itoa(limit))

	return r.URL.Path + "?" + q.Encode()
}

