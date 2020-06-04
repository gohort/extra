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

	// When using a typed extra map, you should keep in mind that the other
	// extra fields that do not match the type will be omitted by default.
	// If you wish to keep the rest look at the `multiple-maps` example.
	extraInts extra.Ints // only captures the integers.
}

// UnmarshalJSON is an override to the json.Unmarshal function and will use what
// you place inside this definition instead.
func (a *Info) UnmarshalJSON(data []byte) error {
	// unmarshal the json using the `extra` package unmarshal function.
	// When unmarshalling pass in the data, structure, and the extra maps.
	return extra.Unmarshal(data, a, &a.extraInts)
}

// MarshalJSON is an override to the json.Marshal function and will use what
// you place inside this definition instead.
func (a *Info) MarshalJSON() ([]byte, error) {
	// marshal the structure using the `extra` package marshal function.
	// Pass in the stucture itself and all the maps you wish you expose.
	return extra.Marshal(a, &a.extraInts)
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

	bb, err := json.MarshalIndent(&myStruct, "", "    ")
	if err != nil {
		panic(err)
	}

	// Print out the JSON showing that it is missing the extra string field.
	//
	// The JSON is missing the extra string value, because there is no way that
	// the integers map can store the value.
	//
	// To capture it we could either use a `extra.Any` map
	// or could add another map to capture strings.
	fmt.Println()
	fmt.Printf("%s", bb)
}
