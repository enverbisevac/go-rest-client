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

Library provide three basic free functions: Modify, Get and Delete. If you want to retrieve data just use
```go
type Article struct {
   Title string `json:"title"`
   Body  string `json:"body"`
}

article, err := Get[Article](context.Background(), "your resource url")
```

For creating or modifying resource:
```go
type Article struct {
   Title string `json:"title"`
   Body  string `json:"body"`
}

article := Article{
	Title: "some title",
	Body: "some body"
}

res, err := Modify[Article](context.Background(), http.MethodPost, "your resource url", rest.WithBody(&article))
```
