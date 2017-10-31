package pagination

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// Search is a pagination search structure.
type Search struct {
	SQL        string
	Parameters []interface{}
}

// NewSearch uses the url parameters to create a search struct.
// Returns an unknown search operation if an unknown search operation was encountered.
// Returns a search operator is missing error if multiple search condition were provided without a search operation.
// Returns a cannot find search conditions error if a search operator was provided without having at least two search conditions.
func NewSearch(query url.Values) (*Search, error) {
	var conditions []string
	var parameters []interface{}

	operator := query.Get("searchOperator")

	for queryParam, value := range query {
		if isASearchCondition, _ := regexp.MatchString(`^(.+__.+)$`, queryParam); isASearchCondition && len(value) != 0 {
			paramComponents := strings.Split(queryParam, "__")
			condition, parameter := GetSearchComponents(paramComponents[0], paramComponents[1], value[0])

			if condition == "" || parameter == "" {
				return nil, errors.New("Unknown search operation '" + paramComponents[1] + "'")
			}

			conditions = append(conditions, condition)
			parameters = append(parameters, parameter)
		}
	}

	if len(conditions) > 1 && operator == "" {
		return nil, errors.New("Search operator is missing")
	}

	if operator != "" && len(conditions) < 2 {
		return nil, errors.New("Cannot find search conditions")
	}

	// No search query parameters were provided in the url
	if len(conditions) == 0 && len(parameters) == 0 {
		return nil, nil
	}

	return &Search{
		SQL:        "(" + strings.Join(conditions, " "+operator+" ") + ")",
		Parameters: parameters,
	}, nil
}

// GetSearchComponents is a helper method that returns a search condition.
func GetSearchComponents(field, operator, value string) (condition, parameter string) {
	switch operator {
	case "equals":
		condition = "(" + field + " = ?)"
		parameter = value
	case "notequals":
		condition = "(" + field + " != ?)"
		parameter = value
	case "greaterthan":
		condition = "(" + field + " > ?)"
		parameter = value
	case "lessthan":
		condition = "(" + field + " < ?)"
		parameter = value
	case "gthanorequals":
		condition = "(" + field + " >= ?)"
		parameter = value
	case "lthanorequals":
		condition = "(" + field + " <= ?)"
		parameter = value
	case "startswith":
		condition = "(" + field + " LIKE ?)"
		parameter = value + "%"
	case "endswith":
		condition = "(" + field + " LIKE ?)"
		parameter = "%" + value
	case "contains":
		condition = "(" + field + " LIKE ?)"
		parameter = "%" + value + "%"
	case "after":
		condition = "(" + field + " > ?)"
		parameter = value
	case "before":
		condition = "(" + field + " < ?)"
		parameter = value
	case "year":
		condition = "(YEAR(" + field + ") = ?)"
		parameter = value
	case "month":
		condition = "(MONTH(" + field + ") = ?)"
		parameter = value
	case "day":
		condition = "(DAY(" + field + ") = ?)"
		parameter = value
	default:
		condition = ""
		parameter = ""
	}

	return condition, parameter
}
