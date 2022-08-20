package main

import (
	"context"
	"github.com/enverbisevac/go-rest-client"
	"log"
	"time"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

type Response[T any] struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
	Data       []T `json:"data"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, err := rest.Get[Response[User]](ctx, "https://reqres.in/api/users")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %+v", response)
}
