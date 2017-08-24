package pagination

import (
	"net/url"
	"strconv"
)

// Links is a pagination Links structure.
type Links struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Self     string `json:"self"`
}

// NewLinks creates pagination links.
func NewLinks(reqURL *url.URL, count int) *Links {
	query := reqURL.Query()
	Links := &Links{Self: reqURL.String()}
	page, err := strconv.ParseInt(query.Get("page"), 10, 64)
	// A page number is given
	if err == nil {
		// Next Links
		limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
		if err == nil && int64(count) >= limit {
			query.Set("page", strconv.Itoa(int(page+1)))
			reqURL.RawQuery = query.Encode()
			Links.Next = reqURL.String()
		}
		// Previous Links
		if page > 1 {
			query.Set("page", strconv.Itoa(int(page-1)))
			reqURL.RawQuery = query.Encode()
			Links.Previous = reqURL.String()
		}
	}

	return Links
}
