// Package page provides support for query paging.
package page

// Document is the form used for API responses from query API calls.
type Document[T any] struct {
	Items       []T `json:"items"`
	Total       int `json:"total"`
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
}

// NewDocument constructs a response value for a web paging trusted.
func NewDocument[T any](items []T, total int, page int, rowsPerPage int) Document[T] {
	return Document[T]{
		Items:       items,
		Total:       total,
		Page:        page,
		RowsPerPage: rowsPerPage,
	}
}
