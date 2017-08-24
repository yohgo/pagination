package pagination_test

import (
	"net/url"
	"reflect"
	"testing"

	"bitbucket.org/effcommsa/illuminate-common/pagination"
)

// newLinksDataProvider provides data for the TestNewLinks function.
var newLinksDataProvider = []struct {
	name  string
	url   string
	count int
	links *pagination.Links
}{
	{
		name:  "Successful links creation - no paging, no ordering",
		url:   "api.demo.com/v1/users",
		count: 3,
		links: &pagination.Links{
			Next:     "",
			Previous: "",
			Self:     "api.demo.com/v1/users",
		},
	},
	{
		name:  "Successful links creation - ordering, no paging",
		url:   "api.demo.com/v1/users?order_by=name&order=asc",
		count: 3,
		links: &pagination.Links{
			Next:     "",
			Previous: "",
			Self:     "api.demo.com/v1/users?order_by=name&order=asc",
		},
	},
	{
		name:  "Successful links creation - first page",
		url:   "api.demo.com/v1/users?page=1&limit=3&order_by=name&order=asc",
		count: 3,
		links: &pagination.Links{
			Next:     "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=2",
			Previous: "",
			Self:     "api.demo.com/v1/users?page=1&limit=3&order_by=name&order=asc",
		},
	},
	{
		name:  "Successful links creation - second page",
		url:   "api.demo.com/v1/users?page=2&limit=3&order_by=name&order=asc",
		count: 3,
		links: &pagination.Links{
			Next:     "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=3",
			Previous: "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=1",
			Self:     "api.demo.com/v1/users?page=2&limit=3&order_by=name&order=asc",
		},
	},
	{
		name:  "Successful links creation - third page",
		url:   "api.demo.com/v1/users?page=3&limit=3&order_by=name&order=asc",
		count: 3,
		links: &pagination.Links{
			Next:     "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=4",
			Previous: "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=2",
			Self:     "api.demo.com/v1/users?page=3&limit=3&order_by=name&order=asc",
		},
	},
	{
		name:  "Successful links creation - last page",
		url:   "api.demo.com/v1/users?page=4&limit=3&order_by=name&order=asc",
		count: 2,
		links: &pagination.Links{
			Next:     "",
			Previous: "api.demo.com/v1/users?limit=3&order=asc&order_by=name&page=3",
			Self:     "api.demo.com/v1/users?page=4&limit=3&order_by=name&order=asc",
		},
	},
}

// TestNewLinks tests the paginator NewLinks method.
func TestNewLinks(t *testing.T) {
	t.Log("NewLinks")
	// Check each test case
	for _, testcase := range newLinksDataProvider {
		t.Log(testcase.name)

		url, _ := url.Parse(testcase.url)
		links := pagination.NewLinks(url, testcase.count)

		// Check links
		if !reflect.DeepEqual(testcase.links, links) {
			t.Errorf("Expected links to be %q but got %q", testcase.links, links)
		}
	}
}
