package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"min=1,max=20"`
	Offset int      `json:"offset" validate:"min=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
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

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	} else {
		fq.Tags = []string{}
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	return *fq, nil
}
