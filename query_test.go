package pagination_test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/yohgo/pagination"
)

// newQueryDataProvider provides data for the TestNewQuery function.
var newQueryDataProvider = []struct {
	name  string
	query string
	count int
	got   *pagination.Query
	err   error
}{
	{
		name:  "Query creation fails due to an invalid url query",
		query: "page=invalid_page&limit=invalid_limit&order_by=&order=invalid_order",
		got:   nil,
		err:   errors.New("Page is invalid"),
	},
	{
		name:  "Successful Query creation - no paging, no ordering",
		query: "api.demo.com/v1/users",
		got: &pagination.Query{
			Page:    0,
			Limit:   0,
			OrderBy: "",
			Order:   "",
		},
		err: nil,
	},
	{
		name:  "Successful Query creation - ordering, no paging",
		query: "order_by=name&order=asc",
		got: &pagination.Query{
			Page:    0,
			Limit:   0,
			OrderBy: "name",
			Order:   "asc",
		},
		err: nil,
	},
	{
		name:  "Successful Query creation - first page",
		query: "page=1&limit=3&order_by=name&order=asc",
		got: &pagination.Query{
			Page:    1,
			Limit:   3,
			OrderBy: "name",
			Order:   "asc",
		},
		err: nil,
	},
	{
		name:  "Successful Query creation - second page",
		query: "page=2&limit=3&order_by=name&order=asc",
		got: &pagination.Query{
			Page:    2,
			Limit:   3,
			OrderBy: "name",
			Order:   "asc",
		},
		err: nil,
	},
	{
		name:  "Successful Query creation - third page",
		query: "page=3&limit=3&order_by=surname&order=desc",
		got: &pagination.Query{
			Page:    3,
			Limit:   3,
			OrderBy: "surname",
			Order:   "desc",
		},
		err: nil,
	},
}

// TestNewQuery tests the paginator NewQuery method.
func TestNewQuery(t *testing.T) {
	t.Log("NewQuery")
	// Check each test case
	for _, testcase := range newQueryDataProvider {
		t.Log(testcase.name)

		query, _ := url.ParseQuery(testcase.query)
		want, err := pagination.NewQuery(query)

		// Check error
		if !reflect.DeepEqual(testcase.err, err) {
			t.Errorf("Expected error to be %q but got %q", testcase.err, err)
		}

		// Check query
		if !reflect.DeepEqual(testcase.got, want) {
			t.Errorf("Expected query to be %q but got %q", testcase.got, want)
		}
	}
}

// validateQueryDataProvider provides data for the TestValidateQuery function.
var validateQueryDataProvider = []struct {
	name  string
	query string
	err   error
}{
	{
		name:  "A successful validation",
		query: "page=1&limit=5&order_by=name&order=desc",
		err:   nil,
	},
	{
		name:  "A successful validation with no ordering",
		query: "page=1&limit=5&order_by=&order=",
		err:   nil,
	},
	{
		name:  "Validation fails due to an invalid page",
		query: "page=invalid page&limit=invalid limit&order_by=&order=",
		err:   errors.New("Page is invalid"),
	},
	{
		name:  "Validation fails due to a zero page",
		query: "page=0&limit=invalid limit&order_by=&order=",
		err:   errors.New("Page is invalid"),
	},
	{
		name:  "Validation fails due to a negative page",
		query: "page=-3&limit=invalid limit&order_by=&order=",
		err:   errors.New("Page is invalid"),
	},
	{
		name:  "Validation fails due to a missing page",
		query: "page=&limit=2&order_by=&order=",
		err:   errors.New("Page is missing"),
	},
	{
		name:  "Validation fails due to an invalid limit",
		query: "page=1&limit=invalid limit&order_by=&order=",
		err:   errors.New("Limit is invalid"),
	},
	{
		name:  "Validation fails due to an zero limit",
		query: "page=1&limit=0&order_by=&order=",
		err:   errors.New("Limit is invalid"),
	},
	{
		name:  "Validation fails due to a negative limit",
		query: "page=1&limit=-4&order_by=&order=",
		err:   errors.New("Limit is invalid"),
	},
	{
		name:  "Validation fails due to a missing order by",
		query: "page=1&limit=5&order_by=&order=desc",
		err:   errors.New("Order by is missing"),
	},
	{
		name:  "Validation fails due to an invalid order",
		query: "page=1&limit=5&order_by=name&order=invalid order",
		err:   errors.New("Order is invalid"),
	},
}

// TestValidateQuery tests the paginator ValidateQuery method.
func TestValidateQuery(t *testing.T) {
	t.Log("ValidateQuery")
	// Check each test case
	for _, testcase := range validateQueryDataProvider {
		t.Log(testcase.name)

		query, _ := url.ParseQuery(testcase.query)
		err := pagination.ValidateQuery(query)

		// Check error
		if !reflect.DeepEqual(testcase.err, err) {
			t.Errorf("Expected error to be %q but got %q", testcase.err, err)
		}
	}
}

// getOrderDataProvider provides data for the TestGetOrder function.
var getOrderDataProvider = []struct {
	name  string
	query *pagination.Query
	want  string
}{
	{
		name: "An order retrieval with complete query",
		query: &pagination.Query{
			Page:    1,
			Limit:   5,
			OrderBy: "name",
			Order:   "desc",
		},
		want: "name desc",
	},
	{
		name: "An order retrieval with query having no order",
		query: &pagination.Query{
			Page:    0,
			Limit:   0,
			OrderBy: "name",
		},
		want: "name asc",
	},
	{
		name: "An order retrieval with query having no order by",
		query: &pagination.Query{
			Page:  3,
			Limit: 6,
			Order: "desc",
		},
		want: "created_at desc",
	},
	{
		name:  "An order retrieval with query having no order and order by",
		query: &pagination.Query{},
		want:  "created_at asc",
	},
}

// TestGetOrder tests the paginator GetOrder method.
func TestGetOrder(t *testing.T) {
	t.Log("GetOrder")
	// Check each test case
	for _, testcase := range getOrderDataProvider {
		t.Log(testcase.name)

		got := testcase.query.GetOrder()

		// Check response
		if !reflect.DeepEqual(testcase.want, got) {
			t.Errorf("Expected response to be %s but got %s", testcase.want, got)
		}
	}
}

// getLimitDataProvider provides data for the TestGetLimit function.
var getLimitDataProvider = []struct {
	name  string
	query *pagination.Query
	want  int
}{
	{
		name: "A limit retrieval with complete query",
		query: &pagination.Query{
			Page:    1,
			Limit:   5,
			OrderBy: "name",
			Order:   "desc",
		},
		want: 5,
	},
	{
		name: "A limit retrieval with query having no limit",
		query: &pagination.Query{
			Page:    0,
			OrderBy: "name",
			Order:   "",
		},
		want: 30,
	},
	{
		name: "A limit retrieval with query having 1 as limit",
		query: &pagination.Query{
			Page:    0,
			Limit:   1,
			OrderBy: "",
			Order:   "desc",
		},
		want: 1,
	},
	{
		name: "A limit retrieval with query having 1000 as limit",
		query: &pagination.Query{
			Page:    0,
			Limit:   1000,
			OrderBy: "",
			Order:   "desc",
		},
		want: 1000,
	},
	{
		name: "A limit retrieval with query having a zero limit",
		query: &pagination.Query{
			Page:    0,
			Limit:   0,
			OrderBy: "",
			Order:   "desc",
		},
		want: 30,
	},
	{
		name: "A limit retrieval with query having a negative limit",
		query: &pagination.Query{
			Page:    0,
			Limit:   0,
			OrderBy: "",
			Order:   "desc",
		},
		want: 30,
	},
}

// TestGetLimit tests the paginator GetLimit method.
func TestGetLimit(t *testing.T) {
	t.Log("GetLimit")
	// Check each test case
	for _, testcase := range getLimitDataProvider {
		t.Log(testcase.name)

		got := testcase.query.GetLimit()

		// Check response
		if !reflect.DeepEqual(testcase.want, got) {
			t.Errorf("Expected response to be %d but got %d", testcase.want, got)
		}
	}
}

// getOffsetDataProvider provides data for the TestGetOffset function.
var getOffsetDataProvider = []struct {
	name  string
	query *pagination.Query
	want  int
}{
	{
		name: "A offset retrieval with complete query",
		query: &pagination.Query{
			Page:    3,
			Limit:   5,
			OrderBy: "name",
			Order:   "desc",
		},
		want: 10,
	},
	{
		name: "A offset retrieval with query having no page",
		query: &pagination.Query{
			Limit:   5,
			OrderBy: "name",
			Order:   "",
		},
		want: 0,
	},
	{
		name: "A offset retrieval with query having 1 as page",
		query: &pagination.Query{
			Page:    1,
			Limit:   10,
			OrderBy: "",
			Order:   "desc",
		},
		want: 0,
	},
	{
		name: "A offset retrieval with query having 2 as page",
		query: &pagination.Query{
			Page:    2,
			Limit:   8,
			OrderBy: "",
			Order:   "desc",
		},
		want: 8,
	},
	{
		name: "A offset retrieval with query having 1000 as page",
		query: &pagination.Query{
			Page:    1000,
			Limit:   3,
			OrderBy: "",
			Order:   "desc",
		},
		want: 2997,
	},
	{
		name: "A offset retrieval with query having a zero page",
		query: &pagination.Query{
			Page:    0,
			Limit:   5,
			OrderBy: "",
			Order:   "desc",
		},
		want: 0,
	},
	{
		name: "A offset retrieval with query having a negative page",
		query: &pagination.Query{
			Page:    -3,
			Limit:   0,
			OrderBy: "",
			Order:   "desc",
		},
		want: 0,
	},
}

// TestGetOffset tests the paginator GetOffset method.
func TestGetOffset(t *testing.T) {
	t.Log("GetOffset")
	// Check each test case
	for _, testcase := range getOffsetDataProvider {
		t.Log(testcase.name)

		got := testcase.query.GetOffset()

		// Check response
		if !reflect.DeepEqual(testcase.want, got) {
			t.Errorf("Expected response to be %d but got %d", testcase.want, got)
		}
	}
}
