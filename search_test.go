package pagination_test

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/yohgo/pagination"
)

// newSearchDataProvider provides data for the TestNewQuery function.
var newSearchDataProvider = []struct {
	name  string
	query string
	err   error
}{
	{
		name:  "Successful search creation",
		query: "name__equals=ammar&type__notequals=admin&age__greaterthan=18&rank__lessthan=3&amount__gthanorequals=1000&level__lthanorequals=5&department__startswith=dev&department__endswith=nt&department__contains=mnt&created__after=2016-08-07&created__before=2016-09-07&created__year=2016&created__month=8&created__day=5&searchOperator=AND",
		err:   nil,
	},
	{
		name:  "Successful search creation - Missing search operator",
		query: "name__equals=ammar&type__notequals=admin&age__greaterthan=18",
		err:   errors.New("Search operator is missing"),
	},
	{
		name:  "Successful search creation - Unknown search operation",
		query: "name__equals=ammar&type__notequals=admin&age__unknownoperation=18&searchOperator=AND",
		err:   errors.New("Unknown search operation 'unknownoperation'"),
	},
	{
		name:  "Successful search creation - No search conditions",
		query: "searchOperator=AND",
		err:   errors.New("Cannot find search conditions"),
	},
	{
		name:  "A failed search creation - No search query parameters",
		query: "page=3&limit=3&order_by=surname&order=desc",
		err:   nil,
	},
}

// TestNewSearch tests the paginator NewSearch method.
func TestNewSearch(t *testing.T) {
	t.Log("NewSearch")
	// Check each test case
	for _, testcase := range newSearchDataProvider {
		t.Log(testcase.name)

		query, _ := url.ParseQuery(testcase.query)
		_, err := pagination.NewSearch(query)

		// Check error
		if !reflect.DeepEqual(testcase.err, err) {
			t.Errorf("Expected error to be %q but got %q", testcase.err, err)
		}
	}
}

// getSearchComponentsDataProvider provides data for the TestgetSearchComponents function.
var getSearchComponentsDataProvider = []struct {
	name      string
	field     string
	operator  string
	value     string
	condition string
	parameter string
}{
	{
		name:      "An search condition retrieval with the equals operation",
		field:     "id",
		operator:  "equals",
		value:     "3",
		condition: "(id = ?)",
		parameter: "3",
	},
	{
		name:      "An search condition retrieval with the notequals operation",
		field:     "id",
		operator:  "notequals",
		value:     "3",
		condition: "(id != ?)",
		parameter: "3",
	},
	{
		name:      "An search condition retrieval with the greaterthan operation",
		field:     "id",
		operator:  "greaterthan",
		value:     "18",
		condition: "(id > ?)",
		parameter: "18",
	},
	{
		name:      "An search condition retrieval with the lessthan operation",
		field:     "id",
		operator:  "lessthan",
		value:     "18",
		condition: "(id < ?)",
		parameter: "18",
	},
	{
		name:      "An search condition retrieval with the gthanorequals operation",
		field:     "id",
		operator:  "gthanorequals",
		value:     "18",
		condition: "(id >= ?)",
		parameter: "18",
	},
	{
		name:      "An search condition retrieval with the lthanorequals operation",
		field:     "id",
		operator:  "lthanorequals",
		value:     "18",
		condition: "(id <= ?)",
		parameter: "18",
	},
	{
		name:      "An search condition retrieval with the startswith operation",
		field:     "name",
		operator:  "startswith",
		value:     "am",
		condition: "(name LIKE ?)",
		parameter: "am%",
	},
	{
		name:      "An search condition retrieval with the endswith operation",
		field:     "name",
		operator:  "endswith",
		value:     "am",
		condition: "(name LIKE ?)",
		parameter: "%am",
	},
	{
		name:      "An search condition retrieval with the contains operation",
		field:     "name",
		operator:  "contains",
		value:     "mm",
		condition: "(name LIKE ?)",
		parameter: "%mm%",
	},
	{
		name:      "An search condition retrieval with the after operation",
		field:     "created_at",
		operator:  "after",
		value:     "2016-08-07 00:00:00",
		condition: "(created_at > ?)",
		parameter: "2016-08-07 00:00:00",
	},
	{
		name:      "An search condition retrieval with the before operation",
		field:     "created_at",
		operator:  "before",
		value:     "2016-08-07 00:00:00",
		condition: "(created_at < ?)",
		parameter: "2016-08-07 00:00:00",
	},
	{
		name:      "An search condition retrieval with the year operation",
		field:     "created_at",
		operator:  "year",
		value:     "2017",
		condition: "(YEAR(created_at) = ?)",
		parameter: "2017",
	},
	{
		name:      "An search condition retrieval with the month operation",
		field:     "created_at",
		operator:  "month",
		value:     "10",
		condition: "(MONTH(created_at) = ?)",
		parameter: "10",
	},
	{
		name:      "An search condition retrieval with the day operation",
		field:     "created_at",
		operator:  "day",
		value:     "11",
		condition: "(DAY(created_at) = ?)",
		parameter: "11",
	},
	{
		name:      "A failed condition retrieval",
		field:     "",
		operator:  "",
		value:     "",
		condition: "",
	},
}

// TestGetSearchComponents tests the paginator GetSearchComponents method.
func TestGetSearchComponents(t *testing.T) {
	t.Log("GetSearchComponents")
	// Check each test case
	for _, testcase := range getSearchComponentsDataProvider {
		t.Log(testcase.name)

		condition, parameter := pagination.GetSearchComponents(testcase.field, testcase.operator, testcase.value)

		// Check condition
		if !reflect.DeepEqual(testcase.condition, condition) {
			t.Errorf("Expected response to be %s but got %s", testcase.condition, condition)
		}

		// Check parameter
		if !reflect.DeepEqual(testcase.parameter, parameter) {
			t.Errorf("Expected parameter to be %s but got %s", testcase.parameter, parameter)
		}
	}
}
