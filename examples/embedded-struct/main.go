package main

import (
	"encoding/json"
	"fmt"

	"github.com/gohort/extra"
)

var myJSON = []byte(`{
	"firstName": "John",
	"lastName": "Doe",
	"age": 20,
	"email": "johndoe@email.com",
	"pass": "***"
}`)

// Person stores structure information about a person.
type Person struct {
	First string `json:"firstName"`
	Last  string `json:"lastName"`
	Age   int    `json:"age"`
}

// Profile is used for example purposes to show using a embedded structures
// with an extras map.
type Profile struct {
	// This field should remain empty since we define it in Person.
	First string `json:"firstName"`
	Pass  string `json:"pass"`
	// Person catches all the fields that are defined in the structure.
	Person
	// extras captures any strings in the JSON that is not currently defined
	// in the golang structure.
	extras extra.Strings
}

// UnmarshalJSON is an override to the json.Unmarshal function and will use what
// you place inside this definition instead.
func (a *Profile) UnmarshalJSON(data []byte) error {
	// unmarshal the json using the `extra` package unmarshal function.
	// When unmarshalling pass in the data, structure, and the extra maps.
	return extra.Unmarshal(data, a, &a.extras)
}

// MarshalJSON is an override to the json.Marshal function and will use what
// you place inside this definition instead.
func (a *Profile) MarshalJSON() ([]byte, error) {
	// marshal the structure using the `extra` package marshal function.
	// Pass in the stucture itself and all the maps you wish you expose.
	return extra.Marshal(a, &a.extras)
}

func main() {
	var profile Profile
	if err := json.Unmarshal(myJSON, &profile); err != nil {
		panic(err)
	}
	// Print that the structure is filled out properly.
	fmt.Printf("%#v\n", profile)

	// Print out using the maps
	fmt.Printf("my name is %s %s\n", profile.Person.First, profile.Person.Last)
	fmt.Printf("my email is %s\n", profile.extras.Get("email"))

	bb, err := json.MarshalIndent(&profile, "", "    ")
	if err != nil {
		panic(err)
	}

	// Print out the JSON being the same as how we got it.
	fmt.Println()
	fmt.Printf("%s", bb)
}
