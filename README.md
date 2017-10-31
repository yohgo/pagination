pagination
=======
[![Build Status](https://travis-ci.org/yohgo/pagination.svg?branch=master)](https://travis-ci.org/yohgo/pagination)
[![goreportcard for yohgo/pagination](https://goreportcard.com/badge/github.com/yohgo/pagination)](https://goreportcard.com/report/yohgo/pagination)
[![codecov for yohgo/pagination](https://codecov.io/gh/yohgo/pagination/branch/master/graph/badge.svg)](https://codecov.io/gh/yohgo/pagination)

yohgo/pagination is a general purpose pagination library ideal for paginating and filtering RESTful response data.

---------------------------------------

  * [Requirements](#requirements)
  * [Features](#features)
  * [Installation](#installation)
  * [Usage](#usage)
    * [Create a Pagination Query](#create-a-query)
    * [Create a Pagination Result Page](#create-a-result-page)

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

Once pagination is [installed](#installation), the following features are available:

### Create a Pagination Query

Creating a pagination query is as simple as calling the `NewQuery` method. The `NewQuery` method take as parameter of type [url.Values](https://golang.org/pkg/net/url/#Values) and returns a pagination `Query` as shown below.

```go
func UsersHandler(w http.ResponseWriter, req *http.Request) {
    query, err := pagination.NewQuery(req.URL.Query())

    if err != nil {
        /* Error handling logic */
    }

    // Passing the query down the layer-stack
    results := UsersService.GetAll(query)
}
```

For example, in the above snippet if `req` has the following raw url `api.awesome.com/users?page=1&limit=30&order_by=id&order=asc&name_contains=dav&age_greaterthan=20&searchOperator=OR` pagination will produce the following `Query` 

```go
query := &pagination.Query{
	Page:    1,
	Limit:   30,
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

The above query struct can be used to filter results at the repository layer. For example, the following snippet uses pagination `Query` and [GORM](http://jinzhu.me/gorm/) to retrieve a paginated slice of users:

```go
func GetAllUsers(query Query) []User {
    var users []User

    // Retrieving a filtered slice of user using GORM
    if err := repository.DB.
        Order(query.GetOrder()).
        Limit(query.GetLimit()).
        Offset(query.GetOffset()).
        Find(&users).Error; err != nil {

        return nil
    }

    return users
}
```

### Create a Pagination Result Page
After using the pagination `Query` to retrieve a result collection; the `NewPage` method can be used to create a pagination `Page` as shown in the following snippet.

```go
func UsersHandler(w http.ResponseWriter, req *http.Request) {
    query, err := pagination.NewQuery(req.URL.Query())

    if err != nil {
        /* Error handling logic */
    }

    // Passing the query to the users service
    users := UsersService.GetAll(query)

    // Now converting the raw result into a pagination page 
    page, _ := pagination.NewPage(req.URL, users)
}
```

For example, in the above snippet if `req` has the following raw url `api.awesome.com/users?page=1&limit=2&order_by=id&order=asc&name_contains=dav&age_greaterthan=20&searchOperator=OR` pagination will produce the following `page`

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
            Name: "user#1",
        },
        {
            Name: "User#2",
        },
    }
}
```
