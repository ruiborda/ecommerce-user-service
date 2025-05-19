package dto

import "strconv"

type Pageable struct {
	Page  int    `json:"page"`
	Size  int    `json:"size"`
	Query string `json:"query"`
}

func NewPageable(pageStr string, sizeStr string, query string) *Pageable {
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	return &Pageable{Page: page, Size: size, Query: query}
}
