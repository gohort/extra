package main

import (
	"fmt"

	extra "github.com/gohort/extra/v2"
)

type User struct {
	Username string `json:"username"`

	X extra.Map `json:"-"`
}

var data = []byte(`{
	"username": "user123",
	"email": "user123@email.com"
}`)

func main() {
	var user User

	if err := extra.Unmarshal(data, &user); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", user)
}