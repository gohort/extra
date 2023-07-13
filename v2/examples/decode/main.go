package main

import (
	"encoding/json"
	"fmt"

	"github.com/gohort/extra/v2"
)

type Topic struct {
	Id int `json:"id"`
	Slug string `json:"slug"`
}

type UserTopicsList struct {
	Topics []Topic `json:"topics"`
	MoreTopicsURL string `json:"moreTopicsURL"`
}

type User struct {
	Username string `json:"username"`
	Topics *UserTopicsList `json:"topics"`

	X extra.Map `json:"-"`
}


var data = []byte(`{
	"username": "testUsername",
	"topics": {
		"moreTopicsURL": "http://example.com/more",
		"topics": [
			{ "id": 123456789, "slug": "salt1" },
			{ "id": 123456789, "slug": "salt2" },
			{ "id": 123456789, "slug": "salt3" },
			{ "id": 123456789, "slug": "salt4" },
			{ "id": 123456789, "slug": "salt5" },
			{ "id": 123456789, "slug": "salt6" }
		]
	},
	"extra": "thing",
	"here": 123
}`)

func main() {
	var user User

	if err := json.Unmarshal(data, &user.X); err != nil {
		panic(err)
	}

	// if err := extra.Decode(&user.X, &user); err != nil {
	// 	panic(err)
	// }

	fmt.Printf("%#v\n", user)
}