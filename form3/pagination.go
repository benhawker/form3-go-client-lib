package form3

import (
	"net/url"
	"strconv"
)

const (
	defaultPageNumber int = 0
	defaultSize       int = 10
)

// Pagination -> pagination!
type Pagination struct {
	Number int
	Size   int
}

// NewPagination -> Creates a new Pagination struct
func NewPagination() Pagination {
	return Pagination{defaultPageNumber, defaultSize}
}

// Params -> sets pagination values
func (p *Pagination) Params() url.Values {
	params := url.Values{}

	number := p.Number
	if number < 0 {
		number = 0
	}

	size := p.Size
	if size <= 0 {
		size = 100
	}

	params.Add("page[number]", strconv.Itoa(number))
	params.Add("page[size]", strconv.Itoa(size))

	return params
}
