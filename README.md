pagination
=======
[![Build Status](https://travis-ci.org/yohgo/pagination.svg?branch=master)](https://travis-ci.org/yohgo/pagination)
[![goreportcard for yohgo/pagination](https://goreportcard.com/badge/github.com/yohgo/pagination)](https://goreportcard.com/report/yohgo/pagination)
[![codecov for yohgo/pagination](https://codecov.io/gh/yohgo/pagination/branch/master/graph/badge.svg)](https://codecov.io/gh/yohgo/pagination)

yohgo/pagination is a general purpose pagination library ideal for paginating and filtering RESTful response data.

![Pagination](pagination.png)

---------------------------------------

  * [Requirements](#requirements)
  * [Features](#features)
  * [Installation](#installation)
  * [Yohgo Pagination Usage](#usage)
    * [Paginating and Filtering Results Using URL Parameters](#filtering-results-using-url-params)
    * [Creating a Pagination Query (Presentation Layer)](#create-a-query)
    * [Handling a Pagination Query (Data Access Layer)](#handle-a-query)

---------------------------------------

## Requirements
  * Go 1.8+

---------------------------------------

## Features

  * Lightweight and fast
  * Native Go implementation
  * Supports both MSSQL and MySQL databases

---------------------------------------

## Installation

Simply install the package to your [$GOPATH](https://github.com/golang/go/wiki/GOPATH "GOPATH") with the [go tool](https://golang.org/cmd/go/ "go command") from shell:

```bash
$ go get github.com/yohgo/pagination
```

---------------------------------------

## Usage

Once yohgo pagination is [installed](#installation), it can be used following the steps below:

### Paginating and Filtering Results Using URL Parameters

Yohgo Pagination provides a set of special `URL Parameters` for the purpose of paginating and filtering the results of an API `GET` requests.

**Paginating Results**

Yohgo Pagination enables users to easily slice results into a set of fixed size `page`'s as well as ordering these `page`'s based on any field/order combination. Yohgo pagination does so by providing the `page`, `limit`, `order_by`, and `order` parameters described below:

Argument | Type | Required? | Description  |  Example
------------ | ------------- | ------------- | ------------- | -------------
page | integer | No | Specifies the current result page | ?page=1
Limit | integer | Yes, if page is set | Specifies the maximum no of records per page | ?page=1&limit=2
order_by | string | No | Sorts the records list by a particular attribute | ?order_by=name
order | string | Yes, if order_by is set | Specifies the sorting direction | ?order_by=name&order=asc

For example, if we have the `api.awesome.com/users` endpoint that manages users, and we want to get a collection of users divided into `page`'s of size 10 and ordered by user name in an ascending fashion, we can do the following:

```
curl -X GET 'http://api.awesome.com/users?page=1&limit=10&order_by=name&order=asc'
```

**Filtering Results**

Similarly, Yohgo Pagination enables users to easily filter results through special `Search Parameters`. A `Search Parameter` follows the following format:

```
{field}__{operator}={value}
```

Also, multiple `Search Parameter`'s can be combined together using the special `searchOperator` parameter which can either be `AND` or `OR`.

```
?{field1}__{operator1}={value1}&{field2}__{operator2}={value2}&searchOperator=AND
```

The following table shows the list of all possible operators provided by yohgo pagination:

| Operator      | Description                                                                                         | Data-types        |
|---------------|-----------------------------------------------------------------------------------------------------|-------------------|
| equals        | The value of the field is tested for `equality` against the specified value.                        | numerics, strings |
| notequals     | The value of the field is tested for `non-equality` against the specified value.                    | numerics, strings |
| greaterthan   | Checks to see whether the value of the field is `greater than` the specified value.                 | numerics          |
| lessthan      | Checks to see whether the value of the field is `less than` the specified value.                    | numerics          |
| gthanorequals | Checks to see whether the value of the field is `greater than or equal to` the specified value.     | numerics          |
| lthanorequals | Checks to see whether the value of the field is `less than or equal to` the specified value.        | numerics          |
| startswith    | Checks to see whether the value of the field `starts with` the specified value.                     | strings           |
| endswith      | Checks to see whether the value of the field `ends with` the specified value.                       | strings           |
| contains      | Checks to see whether the value of the field `contains` the specified value.                        | strings           |
| after         | Filters so that the results have a field value `after` the specified time and date.                 | dates             |
| before        | Filters so that the results have a field value `before` the specified time and date.                | dates             |
| year          | Filters the results so that the `year` section of the field matches the specified numerical value.  | dates             |
| month         | Filters the results so that the `month` section of the field matches the specified numerical value. | dates             |
| day           | Filters the results so that the `day` section of the field matches the specified numerical value.   | dates             |

For example, if we have the `api.awesome.com/users` endpoint that manages users, and we want to get a collection of users where user name contains the string "dav" and divided into `page`'s of size 10, we simply do the following:

```
curl -X GET 'http://api.awesome.com/users?page=1&limit=10&name__contains=dav'
```

### Creating a Pagination Query (Presentation Layer)

Creating a pagination query is as simple as calling the `NewQuery` method. The `NewQuery` method take as parameter of type [url.Values](https://golang.org/pkg/net/url/#Values) and returns a pagination `Query` as shown the simple handler function below.

```go
func UsersHandler(w http.ResponseWriter, req *http.Request) {
    // First, creating a pagination query for the http request
    query, _ := pagination.NewQuery(req.URL.Query())

    // Second, passing the pagination query down the layer stack
    users := UsersService.GetAll(query)

    // Finally, converting the raw result into a pagination page
    page, _ := pagination.NewPage(req.URL, users)

    // ResolveJSON writes a json encoded page into the response (github.com/yohgo/pastry)
    ResolveJSON(w, http.StatusOK, page)
}
```

For example, in the above snippet if `req` has the following raw url `api.awesome.com/users?page=1&limit=2&order_by=id&order=asc&name_contains=dav&age_greaterthan=20&searchOperator=OR` pagination will produce the following `Query` struct:

```go
query := &pagination.Query{
	Page:    1,
	Limit:   2,
	OrderBy: "id",
	Order:   "asc",
	Search:  &pagination.Search{
        SQL: "((name LIKE ?) OR (age > ?))",
        Parameters: []interface{}{
            "%dav%",
            "20",
        }
    },
}
```

And after receiving the paginated/filtered results, yohgo pagination will use this result to produce the following pagination page `page`

```go
page := &pagination.Page{
	Links   &pagination.Links{
        Next: "api.awesome.com/users?page=2&limit=2&order_by=id&order=asc&name_contains=dav&age_greaterthan=20&searchOperator=OR",
        Previous: "",
        Self: "api.awesome.com/users?page=1&limit=2&order_by=id&order=asc&name_contains=dav&age_greaterthan=20&searchOperator=OR",
    },
	Count:   2,
	Results []User{
        {
            Name: "David",
            Age: 23
        },
        {
            Name: "dave",
            Age: 26
        },
    }
}
```

### Handling a Pagination Query (Data Access Layer)

When received from the layers above, the pagination query can be used at the data access layer to dictate how the data is retrieved form the data source, thus, paginating/filtering the results . For example, the following snippet uses pagination a pagination `Query` and [GORM](http://jinzhu.me/gorm/) to retrieve a paginated/filtered slice of users:

```go
func GetAllUsers(query Query) []User {
    var users []User

    // Retrieving a filtered slice of user using GORM
    if err := repository.DB.
        Where(query.Search.SQL, query.Search.Parameters...).
        Order(query.GetOrder()).
        Limit(query.GetLimit()).
        Offset(query.GetOffset()).
        Find(&users).Error; err != nil {

        return nil

    }

    return users
}
```
