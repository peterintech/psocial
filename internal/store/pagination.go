package store

import (
	"fmt"
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"min=1,max=20"`
	Offset int    `json:"offset" validate:"min=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq *PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			return PaginatedFeedQuery{}, fmt.Errorf("invalid limit value: %v", err)
		}
		fq.Limit = parsedLimit
	}

	offset := qs.Get("offset")
	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			return PaginatedFeedQuery{}, fmt.Errorf("invalid offset value: %v", err)
		}
		fq.Offset = parsedOffset
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	return *fq, nil
}
