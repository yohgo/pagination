package pagination_test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
)

// User contains the data required for a unique user.
type User struct {
	ID      uint64
	Name    string
	Surname string
}

// newPageDataProvider provides data for the TestNewPage function.
var newPageDataProvider = []struct {
	name    string
	url     string
	results interface{}
	want    *pagination.Page
	err     error
}{
	{
		name: "Page creation fails due to an invalid url query",
		url:  "api.demo.com/v1/users?page=invalid_page&limit=invalid_limit&order_by=&order=invalid_order",
		want: nil,
		err:  errors.New("Page is invalid"),
	},
	{
		name: "Successful page creation - no paging, no ordering",
		url:  "api.demo.com/v1/users",
		results: []*User{
			{ID: 1, Name: "John", Surname: "Smith"},
		},
		want: &pagination.Page{
			Count: 1,
			Links: &pagination.Links{
				Self:     "api.demo.com/v1/users",
				Previous: "",
				Next:     "",
			},
			Results: []*User{
				{ID: 1, Name: "John", Surname: "Smith"},
			},
		},
		err: nil,
	},
	{
		name: "Successful page creation - ordering, no paging",
		url:  "api.demo.com/v1/users?order_by=name&order=asc",
		results: []*User{
			{ID: 1, Name: "John", Surname: "Smith"},
			{ID: 2, Name: "Jill", Surname: "Doe"},
			{ID: 3, Name: "Paul", Surname: "Johnson"},
		},
		want: &pagination.Page{
			Count: 3,
			Links: &pagination.Links{
				Self:     "api.demo.com/v1/users?order_by=name&order=asc",
				Previous: "",
				Next:     "",
			},
			Results: []*User{
				{ID: 1, Name: "John", Surname: "Smith"},
				{ID: 2, Name: "Jill", Surname: "Doe"},
				{ID: 3, Name: "Paul", Surname: "Johnson"},
			},
		},
		err: nil,
	},
	{
		name:    "Successful page creation - invalid results",
		url:     "api.demo.com/v1/users?page=1&limit=3&order_by=name&order=asc",
		results: "invalid results",
		want:    nil,
		err:     errors.New("The provided collection is not a slice"),
	},
	{
		name: "Successful page creation - first page",
		url:  "api.demo.com/v1/users?page=1&limit=3&order_by=name&order=asc",
		results: []*User{
			{ID: 1, Name: "John", Surname: "Smith"},
			{ID: 2, Name: "Jill", Surname: "Doe"},
			{ID: 3, Name: "Paul", Surname: "Johnson"},
		},
		want: &pagination.Page{
			Count: 3,
			Links: &pagination.Links{
				Next:     "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=2",
				Previous: "",
				Self:     "api.demo.com/v1/users?page=1&limit=3&order_by=name&order=asc",
			},
			Results: []*User{
				{ID: 1, Name: "John", Surname: "Smith"},
				{ID: 2, Name: "Jill", Surname: "Doe"},
				{ID: 3, Name: "Paul", Surname: "Johnson"},
			},
		},
		err: nil,
	},
	{
		name: "Successful page creation - second page",
		url:  "api.demo.com/v1/users?page=2&limit=3&order_by=name&order=asc",
		results: []*User{
			{ID: 1, Name: "John", Surname: "Smith"},
			{ID: 2, Name: "Jill", Surname: "Doe"},
			{ID: 3, Name: "Paul", Surname: "Johnson"},
		},
		want: &pagination.Page{
			Count: 3,
			Links: &pagination.Links{
				Next:     "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=3",
				Previous: "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=1",
				Self:     "api.demo.com/v1/users?page=2&limit=3&order_by=name&order=asc",
			},
			Results: []*User{
				{ID: 1, Name: "John", Surname: "Smith"},
				{ID: 2, Name: "Jill", Surname: "Doe"},
				{ID: 3, Name: "Paul", Surname: "Johnson"},
			},
		},
		err: nil,
	},
	{
		name: "Successful page creation - third page",
		url:  "api.demo.com/v1/users?page=3&limit=3&order_by=surname&order=desc",
		results: []*User{
			{ID: 1, Name: "John", Surname: "Smith"},
			{ID: 2, Name: "Jill", Surname: "Doe"},
		},
		want: &pagination.Page{
			Count: 2,
			Links: &pagination.Links{
				Next:     "",
				Previous: "api.demo.com/v1/users?limit=3&order=desc&order_by=surname&page=2",
				Self:     "api.demo.com/v1/users?page=3&limit=3&order_by=surname&order=desc",
			},
			Results: []*User{
				{ID: 1, Name: "John", Surname: "Smith"},
				{ID: 2, Name: "Jill", Surname: "Doe"},
			},
		},
		err: nil,
	},
}

// TestNewPage tests the paginator NewPage method.
func TestNewPage(t *testing.T) {
	t.Log("NewPage")
	// Check each test case
	for _, testcase := range newPageDataProvider {
		t.Log(testcase.name)

		url, _ := url.Parse(testcase.url)
		got, err := pagination.NewPage(url, testcase.results)

		// Check error
		if !reflect.DeepEqual(testcase.err, err) {
			t.Errorf("Expected error to be %q but got %q", testcase.err, err)
		}

		// Check page
		if !reflect.DeepEqual(testcase.want, got) {
			t.Errorf("Expected page to be %q but got %q", testcase.want, got)
		}
	}
}
