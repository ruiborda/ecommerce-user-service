package dto

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PaginationResponse[T any] struct {
	Links PageLinks `json:"links"`
	Data  *[]*T     `json:"data"`
	Page  Page      `json:"page"`
}

func NewPaginationResponse[T any](c *gin.Context, data *[]*T, totalElements int, pageable *Pageable) *PaginationResponse[T] {
	totalPages := (totalElements + pageable.Size - 1) / pageable.Size

	baseURL := c.Request.URL.Path
	query := c.Request.URL.Query()
	query.Set("page", strconv.Itoa(pageable.Page))
	query.Set("size", strconv.Itoa(pageable.Size))
	currentURL := baseURL + "?" + query.Encode()

	nextURL := ""
	if pageable.Page < totalPages {
		query.Set("page", strconv.Itoa(pageable.Page+1))
		nextURL = baseURL + "?" + query.Encode()
	}

	prevURL := ""
	if pageable.Page > 1 {
		query.Set("page", strconv.Itoa(pageable.Page-1))
		prevURL = baseURL + "?" + query.Encode()
	}

	return &PaginationResponse[T]{
		Links: PageLinks{
			Self: currentURL,
			Next: nextURL,
			Prev: prevURL,
		},
		Data: data,
		Page: Page{
			CurrentPage:   pageable.Page,
			Size:          pageable.Size,
			TotalElements: totalElements,
			TotalPages:    totalPages,
		},
	}
}

type PageLinks struct {
	Self string `json:"self"`
	Next string `json:"next"`
	Prev string `json:"prev"`
}

type Page struct {
	CurrentPage   int `json:"currentPage"`
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}
