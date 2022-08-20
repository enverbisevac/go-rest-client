[![Coverage Status](https://coveralls.io/repos/github/enverbisevac/go-rest-client/badge.svg)](https://coveralls.io/github/enverbisevac/go-rest-client)
# Go REST client library
Simple Go REST client library using generics. Minimum golang version requirement is 1.18

## Install

Installation can be done with a normal go get:
```
go get github.com/enverbisevac/go-rest-client
```

## Usage

First import dependency:
```go
import "github.com/enverbisevac/go-rest-client"
```

You can access package functions with `rest` or if you like aliased imports then use alias.
Library provide three basic top level functions: Modify, Get and Delete.
```go
func Modify[T any](ctx context.Context, method string, requestURL string, options ...Option) (val T, err error)
func Get[T any](ctx context.Context, requestURL string, options ...Option) (result T, err error)
func Delete(ctx context.Context, requestURL string, options ...Option) (err error)
```

Available options:
```go
func WithBody(body any) Option
func WithHeaders(header http.Header) Option
```

Simple GET request example:
```go
type Article struct {
   Title string `json:"title"`
   Body  string `json:"body"`
}

article, err := rest.Get[Article](context.Background(), "your resource url")
```

For creating or modifying resource (POST, PUT, PATCH):
```go
article := Article{
	Title: "some title",
	Body: "some body"
}

res, err := rest.Modify[Article](context.Background(), http.MethodPost, "your resource url", rest.WithBody(&article))
```

Delete a resource example:
```go
article, err := rest.Delete(context.Background(), "your resource url")
```

if you need to provide custom headers or body for any of the functions (get resource example):
```go
article, err := rest.Get[Article](context.Background(), "your resource url", 
	rest.WithHeaders(http.Header{
	    // custom headers 
	})
)
```

or modify a resource:
```go
res, err := rest.Modify[Article](context.Background(), http.MethodPost, "your resource url", 
	rest.WithBody(&article), rest.WithHeaders(... some header))
```

there is also helper function for creating headers:
```go
func Header(opts ...Func) http.Header
```

and option functions:
```go
func WithHeader(header http.Header) HeaderOption
func WithAuth(auth AuthType, token string) HeaderOption
func WithContent(value ContentType) HeaderOption
func With(name string, value ...string) HeaderOption
```

for example:
```go
res, err := rest.Modify[Article](context.Background(), http.MethodPost, "your resource url", 
	rest.WithBody(&article), rest.WithHeaders(
		rest.Headers(
			rest.WithAuth(rest.Bearer, key),
		),
	),
)
```