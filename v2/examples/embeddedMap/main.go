package main

import (
	"fmt"

	extra "github.com/gohort/extra/v2"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`

	somethingPrivate string

	// Using an embedded field for your extra.Map could potentially yield to
	// quicker processing times for structs with lots of other fields as long
	// as you pass the struct as a reference.
	extra.Map `json:"-"`
}

var data = []byte(`{
	"username": "admin",
	"password": "password123",
	"email": "admin@email.com",
	"statistics": {
		"accountAge": "1y5m2d",
		"loginTimes": 12
	}
}`)

func main() {
	user := User{somethingPrivate: "wow!"}

	fmt.Printf("%s\n", data)

	// extra.Unmarshal(data, user) //! this doesn't take advantage of the embedded Map, because it doesn't implement the extra.Mapper interface.

	// Passing the user as a reference is important in order to fulfill the
	// extra.Mapper interface.
	if err := extra.Unmarshal(data, &user); err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", user)

	out, err := extra.MarshalIndent(user, "", "	")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", out)

}