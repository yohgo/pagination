package pagination

import (
	"errors"
	"net/url"
	"reflect"
)

// Page is a pagination page structure.
type Page struct {
	Links   *Links      `json:"_links"`
	Count   int         `json:"count"`
	Results interface{} `json:"results"`
}

// NewPage creates a new pagination page.
// Returns a validation error if page creation was not successful.
// Returns a pagination page if page creation was successful.
func NewPage(reqURL *url.URL, result interface{}) (*Page, error) {
	if err := ValidateQuery(reqURL.Query()); err != nil {
		return nil, err
	}
	// Check if result is a slice
	aType := reflect.ValueOf(result)
	if aType.Kind() != reflect.Slice {
		return nil, errors.New("The provided collection is not a slice")
	}

	return &Page{Links: NewLinks(reqURL, aType.Len()), Count: aType.Len(), Results: result}, nil
}
