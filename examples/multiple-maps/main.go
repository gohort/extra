package main

import (
	"encoding/json"
	"fmt"

	"github.com/gohort/extra"
)

// Info is used for example purposes to show using multiple extra maps.
type Info struct {
	Msg  string `json:"msg"`
	User string `json:"user"`

	// When a JSON may contain multiple types that you'd like to capture in an
	// extras map, but you don't want to use an Any map. You may use multiple
	// extra maps to capture the values based on their types.
	extraStrings extra.Strings // only captures the strings.
	extraInts    extra.Ints    // only captures the integers.
}

// UnmarshalJSON is an override to the json.Unmarshal function and will use what
// you place inside this definition instead.
func (a *Info) UnmarshalJSON(data []byte) error {
	// unmarshal the json using the `extra` package unmarshal function.
	// When unmarshalling pass in the data, structure, and the extra maps.
	return extra.Unmarshal(data, a, &a.extraInts, &a.extraStrings)
}

// MarshalJSON is an override to the json.Marshal function and will use what
// you place inside this definition instead.
func (a *Info) MarshalJSON() ([]byte, error) {
	// marshal the structure using the `extra` package marshal function.
	// Pass in the stucture itself and all the maps you wish you expose.
	return extra.Marshal(a, &a.extraInts, &a.extraStrings)
}

var (
	myJSON = []byte(`
	{
		"msg": "hello",
		"user": "damien",
		"day": "Wednesday",
		"age": 20
	}`)
)

func main() {
	var myStruct Info
	if err := json.Unmarshal(myJSON, &myStruct); err != nil {
		panic(err)
	}
	// Print that the structure is filled out properly.
	fmt.Printf("%v\n", myStruct)

	// Print out using the maps
	fmt.Printf("I am %d years old\n", myStruct.extraInts["age"])
	fmt.Printf("Today is %s\n", myStruct.extraStrings["day"])

	bb, err := json.MarshalIndent(&myStruct, "", "    ")
	if err != nil {
		panic(err)
	}

	// Print out the JSON being the same as how we got it.
	fmt.Println()
	fmt.Printf("%s", bb)
}
