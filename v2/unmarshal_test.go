package extra

import (
	vanillaJSON "encoding/json"
	"reflect"
	"testing"

	"github.com/goccy/go-json"
	old "github.com/gohort/extra"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
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

	X Map `json:"-"`
}

type UserMapStructure struct {
	Username string `json:"username"`
	Topics *UserTopicsList `json:"topics"`

	X map[string]any `mapstructure:",remain"`
}

type UserEmbedded struct {
	Username string `json:"username"`
	Topics *UserTopicsList `json:"topics"`

	Map `json:"-"`
}

type UserMapTop struct {
	Map `json:"-"`

	Username string `json:"username"`
	Topics *UserTopicsList `json:"topics"`
}

type NonUser struct {
	Username string `json:"username"`
	Topics *UserTopicsList `json:"topics"`
}

type UserOld struct {
	Username string `json:"username"`
	Topics *UserTopicsList `json:"topics"`

	X old.Any
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

func BenchmarkUnmarshalWithMap(b *testing.B) {
	var u User
	for i := 0; i < b.N; i++ {
		err := UnmarshalWithMap(data, &u, &u.X)
		_ = err
	}
}

func BenchmarkUnmarshalV1(b *testing.B) {
	var u UserOld
	for i := 0; i < b.N; i++ {
		err := old.Unmarshal(data, &u, &u.X)
		_ = err
	}
}

func BenchmarkUnmarshalV2(b *testing.B) {
	var u User
	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &u)
		_ = err
	}
}

// func BenchmarkUnmarshalV2Decode(b *testing.B) {
// 	var u User
// 	for i := 0; i < b.N; i++ {
// 		err := json.Unmarshal(data, &u.X)
// 		_ = err
// 		err = Decode(&u.X, &u)
// 		_ = err
// 	}
// }

// func TestUnmarshalV2Decode(t *testing.T) {
// 	var u User
// 	err := json.Unmarshal(data, &u.X)
// 	if err != nil {
// 		t.Fatalf("got %s, expected nil", err)
// 	}
// 	err = Decode(&u.X, &u)
// 	if err != nil {
// 		t.Fatalf("got %s, expected nil", err)
// 	}
// }

func BenchmarkUnmarshalMapStructure(b *testing.B) {
	var u UserMapStructure
	for i := 0; i < b.N; i++ {
		var tmp map[string]any
		err := json.Unmarshal(data, &tmp)
		_ = err
		err = mapstructure.Decode(tmp, &u)
		_ = err
	}
}

func BenchmarkUnmarshalV2Embedded(b *testing.B) {
	var u UserEmbedded
	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &u)
		_ = err
	}
}

func BenchmarkUnmarshalV2MapTop(b *testing.B) {
	var u UserMapTop
	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &u)
		_ = err
	}
}

func BenchmarkUnmarshalV2NonExtra(b *testing.B) {
	var u NonUser
	for i := 0; i < b.N; i++ {
		err := Unmarshal(data, &u)
		_ = err
	}
}

func BenchmarkTagName(b *testing.B) {
	var u User
	f := reflect.ValueOf(u).Type().Field(0)
	for i := 0; i < b.N; i ++ {
		name, ignore := parseTagName(f)
		_ = name
		_ = ignore
	}
}

func BenchmarkUnmarshalStd(b *testing.B) {
	var u User
	for i := 0; i < b.N; i++ {
		err := vanillaJSON.Unmarshal(data, &u)
		_ = err
	}
}

func TestUnmarshalMap(t *testing.T) {
	cases := []struct{
		data []byte
		expectedErr error
		expected map[string]any
	}{
		{
			data: []byte(`{"user":"abc","time":123}`),
			expectedErr: nil,
			expected: map[string]any{
				"user": "abc",
				"time": float64(123), // all numbers become float64 when unmarshalled.
			},
		},
		{
			data: []byte(nil),
			expectedErr: &json.SyntaxError{},
			expected: nil,
		},
	}

	for _, tc := range cases {
		var tmp map[string]any
		err := Unmarshal(tc.data, &tmp)
		if tc.expectedErr != nil {
			assert.ErrorAs(t, err, &tc.expectedErr)
		}
		assert.Len(t, tmp, len(tc.expected))
		for k, v := range tmp {
			val, ok := tc.expected[k]
			assert.True(t, ok)

			assert.Equal(t, val, v)
		}
	}
}

func TestUnmarshalExtra(t *testing.T) {
	// EmbeddedStruct is here to test that we can properly unmarshal an embedded
	// object in a struct.
	type EmbeddedStruct struct { Username string `json:"username"` }
	type result struct {
		EmbeddedStruct `json:"embed"`

		ID int64 `json:"id"`


		// To cover checking if a struct field is unexported.
		_ any

		X Map `json:"-"`

		// This field is below the map to cover a line after parsing the tag
		// and not continuing the loop.
		// if mapper != nil { continue }
		IgnoredField string `json:"-"`

		// This field is below the map to cover a line where if we're still
		// looping through the struct fields and we already have the map then
		// remove this field from the map.
		MoreURL string `json:"moreURL"`
	}

	cases := []struct{
		description string

		data []byte
		expected result
		expectedErr error
		doubleRef bool
	}{
		{
			description: "normal use",
			doubleRef: false,
			data: []byte(`{"id":123,"moreURL":"http://more.url/","extra":"info","time":321321,"embed":{"username":"user"}}`),
			expectedErr: nil,
			expected: result{
				ID: 123,
				MoreURL: "http://more.url/",
				EmbeddedStruct: EmbeddedStruct{
					Username: "user",
				},
				X: Map{
					"extra": "info",
					"time": float64(321321),
				},
			},
		},
		{
			description: "double reference",
			doubleRef: true,
			data: []byte(`{"id":123,"moreURL":"http://more.url/","extra":"info","time":321321,"embed":{"username":"user"}}`),
			expectedErr: nil,
			expected: result{
				ID: 123,
				MoreURL: "http://more.url/",
				EmbeddedStruct: EmbeddedStruct{
					Username: "user",
				},
				X: Map{
					"extra": "info",
					"time": float64(321321),
				},
			},
		},
		{
			description: "error on first unmarshalling struct",
			doubleRef: false,
			data: []byte{},
			expectedErr: &json.SyntaxError{},
			expected: result{},
		},
	}

	for _, tc := range cases {
		var r result
		var err error

		// If we want to test double referencing a struct then create a copy of
		// a reference and then reference the reference.
		if tc.doubleRef {
			rr := &r
			err = Unmarshal(tc.data, &rr)
		} else {
			err = Unmarshal(tc.data, &r)
		}

		if tc.expectedErr != nil {
			assert.ErrorAs(t, err, &tc.expectedErr, tc.description)
		}

		assert.Equal(t, tc.expected.ID, r.ID, tc.description)
		assert.Equal(t, tc.expected.MoreURL, r.MoreURL, tc.description)
		assert.Equal(t, tc.expected.EmbeddedStruct.Username, r.EmbeddedStruct.Username, tc.description)
		assert.Len(t, r.X, len(tc.expected.X), tc.description)

		for k, v := range r.X {
			val, ok := tc.expected.X[k]

			assert.True(t, ok, tc.description)
			assert.Equal(t, val, v, tc.description)
		}
	}
}

func TestUnmarshalNil(t *testing.T) {
	err := Unmarshal(data, nil)
	expected := &ErrInvalidUnmarshal{}
	assert.ErrorAs(t, err, &expected)
}

func TestUnmarshalEmbedded(t *testing.T) {
	data := []byte(`{"status":"OK","extra":123}`)
	type result struct {
		Status string `json:"status"`
		Map `json:"-"`
	}
	var r result

	err := Unmarshal(data, &r)
	assert.Nil(t, err)
	assert.Equal(t, r.Status, "OK")
	assert.Len(t, r.Map, 1)

	val, ok := r.Map["extra"]
	assert.True(t, ok)
	assert.Equal(t, val, float64(123))
}

func TestUnmarshalWithMap(t *testing.T) {
	// EmbeddedStruct is here to test that we can properly unmarshal an embedded
	// object in a struct.
	type EmbeddedStruct struct { Username string `json:"username"` }
	type result struct {
		EmbeddedStruct `json:"embed"`

		ID int64 `json:"id"`


		// To cover checking if a struct field is unexported.
		_ any

		// This field is below the map to cover a line after parsing the tag
		// and not continuing the loop.
		// if mapper != nil { continue }
		IgnoredField string `json:"-"`

		// This field is below the map to cover a line where if we're still
		// looping through the struct fields and we already have the map then
		// remove this field from the map.
		MoreURL string `json:"moreURL"`
	}

	cases := []struct{
		description string

		data []byte
		expected result
		expectedMap *Map
		givenMap *Map
		expectedErr error
		doubleRef bool
	}{
		{
			description: "normal use",
			doubleRef: false,
			data: []byte(`{"id":123,"moreURL":"http://more.url/","extra":"info","time":321321,"embed":{"username":"user"}}`),
			expectedErr: nil,
			givenMap: &Map{},
			expected: result{
				ID: 123,
				MoreURL: "http://more.url/",
				EmbeddedStruct: EmbeddedStruct{
					Username: "user",
				},
			},
			expectedMap: &Map{
				"extra": "info",
				"time": float64(321321),
			},
		},
		{
			description: "double reference",
			doubleRef: true,
			data: []byte(`{"id":123,"moreURL":"http://more.url/","extra":"info","time":321321,"embed":{"username":"user"}}`),
			expectedErr: nil,
			givenMap: &Map{},
			expected: result{
				ID: 123,
				MoreURL: "http://more.url/",
				EmbeddedStruct: EmbeddedStruct{
					Username: "user",
				},
			},
			expectedMap: &Map{
				"extra": "info",
				"time": float64(321321),
			},
		},
		{
			description: "error on first unmarshalling map",
			doubleRef: false,
			data: []byte{},
			expectedErr: &json.SyntaxError{},
			expected: result{},
			givenMap: &Map{},
			expectedMap: &Map{},
		},
		{
			description: "nil map",
			doubleRef: false,
			data: []byte{},
			expectedErr: &json.SyntaxError{},
			expected: result{},
		},
	}

	for _, tc := range cases {
		var r result
		var err error

		// If we want to test double referencing a struct then create a copy of
		// a reference and then reference the reference.
		if tc.doubleRef {
			rr := &r
			err = UnmarshalWithMap(tc.data, &rr, tc.givenMap)
		} else {
			err = UnmarshalWithMap(tc.data, &r, tc.givenMap)
		}

		if tc.expectedErr != nil {
			assert.ErrorAs(t, err, &tc.expectedErr, tc.description)
		}
		assert.Equal(t, tc.expected.ID, r.ID, tc.description)
		assert.Equal(t, tc.expected.MoreURL, r.MoreURL, tc.description)

		if tc.expectedMap == nil {
			assert.Nil(t, tc.givenMap)
			continue
		}

		assert.Equal(t, tc.expected.EmbeddedStruct.Username, r.EmbeddedStruct.Username, tc.description)
		assert.Len(t, *tc.givenMap, len(*tc.expectedMap), tc.description)

		for k, v := range *tc.givenMap {
			val, ok := (*tc.expectedMap)[k]

			assert.True(t, ok, tc.description)
			assert.Equal(t, val, v, tc.description)
		}
	}
}

func TestUnmarshalWithMapNonStruct(t *testing.T) {
	var given map[string]any
	var x Map

	err := UnmarshalWithMap([]byte{}, &given, &x)
	assert.ErrorAs(t, err, &ErrNilMap)
}

func TestUnmarshalWithMapNilStruct(t *testing.T) {
	var x Map
	expectedErr := &ErrInvalidUnmarshal{}

	err := UnmarshalWithMap([]byte{}, nil, &x)

	assert.ErrorAs(t, err, &expectedErr)
}