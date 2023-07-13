// Package extra is made to replicate Rust's serde flatten capabilities in Golang.
//
// The goal is to make a fast and easy way of storing unknown variables of a JSON
// object into a map that can be easily reversed back to its original format.
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/gohort/extra/v2"
//	)
//
//	type User struct {
//		Username string `json:"username"`
//
//		X extra.Map `json:"-"`
//	}
//
//	var data = []byte(`{"username":"user123", "email":"user123@email.com"}`)
//
//	func main() {
//		var user User
//		if err := extra.Unmarshal(data, &user); err != nil {
//			// handle err
//		}
//
//		fmt.Println(user) // {user123 map[email:user123@email.com]}
//	}
package extra

import (
	"errors"
	"reflect"

	"github.com/gohort/extra"
)

var (
	ErrNilInterface = extra.ErrNilInterface
	ErrNilMap = errors.New("nil map provided")
)

type ErrInvalidUnmarshal struct {
	Type reflect.Type
}

func (e *ErrInvalidUnmarshal) Error() string {
	if e.Type == nil {
		return "json: passed nil"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "json: passed non-pointer '" + e.Type.String() + "'"
	}

	return "json: nil '" + e.Type.String() + "'"
}
