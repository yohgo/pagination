package pagination

import (
	"errors"
	"net/url"
	"strconv"
)

// Query is a pagination query structure.
type Query struct {
	Page    int
	Limit   int
	OrderBy string
	Order   string
}

// NewQuery creates a new pagination query.
// Returns a validation error if query creation was not successful.
// Returns a pagination query if page creation was successful.
func NewQuery(query url.Values) (*Query, error) {
	if err := ValidateQuery(query); err != nil {
		return nil, err
	}
	// Converts we validated before so we can ignore errors
	page, _ := strconv.Atoi(query.Get("page"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	return &Query{Page: page, Limit: limit, OrderBy: query.Get("order_by"), Order: query.Get("order")}, nil
}

// ValidateQuery validates a collection of url parameters that form a query.
// Returns a page is invalid error if the page is less than 1 or not an integer.
// Returns a limit is invalid error if the limit is less than 1 or not an integer.
// Returns a order by is invalid error if the order by is not a string.
// Returns a order is invalid error if the order is either "asc", or "desc" or empty.
func ValidateQuery(query url.Values) error {
	page, err := strconv.Atoi(query.Get("page"))

	if query.Get("page") != "" && (err != nil || page <= 0) {
		return errors.New("Page is invalid")
	}

	if query.Get("limit") != "" && query.Get("page") == "" {
		return errors.New("Page is missing")
	}

	limit, err := strconv.Atoi(query.Get("limit"))

	if query.Get("limit") != "" && (err != nil || limit <= 0) {
		return errors.New("Limit is invalid")
	}

	if query.Get("order") != "" && query.Get("order_by") == "" {
		return errors.New("Order by is missing")
	}

	if query.Get("order") != "" && query.Get("order") != "asc" && query.Get("order") != "desc" {
		return errors.New("Order is invalid")
	}

	return nil
}

// GetOrder returns a sensible order by string to use when querying data from a datastore.
// Returns a default query order by of created_at when not requesting a particular order by.
// Returns a default query order of ascending when not requesting a particular order.
func (query *Query) GetOrder() string {
	orderBy := map[bool]string{true: "created_at", false: query.OrderBy}[query.OrderBy == ""]
	order := map[bool]string{true: "asc", false: "desc"}[query.Order != "desc"]

	return orderBy + " " + order
}

// GetLimit returns a sensible limit to use when querying data from a datastore.
// Returns a default query limit when requesting for less than 1 record.
// Returns a set query limit when requesting for a number of records within bounds.
func (query *Query) GetLimit() int {
	if query.Limit < 1 {
		return 30
	}

	return query.Limit
}

// GetOffset returns a sensible offset to use when querying data from a datastore.
// Returns a default query offset when requesting for less than the second page.
// Returns a determined query offset when requesting for more than the second page
func (query *Query) GetOffset() int {
	if query.Page < 2 {
		return 0
	}

	return (query.Page - 1) * query.GetLimit()
}
