package extra

import (
	"reflect"
	"strings"
)

type tagOptions string

// Credit: The Go Authors @ "encoding/json"
// parseTag splits a struct field's quirk tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func gatherTags(d any) []reflectSettings {
	elem := reflect.ValueOf(d).Elem()
	numFields := elem.NumField()

	var (
		tag  string // stores the name of the field in Dgraph.
		tags []reflectSettings
	)

	// loop through elements of struct.
	for i := 0; i < numFields; i++ {
		tag, _ = parseTag(reflect.TypeOf(d).Elem().Field(i).Tag.Get("json"))
		if tag == "" {
			continue
		}

		tags = append(tags, reflectSettings{i, tag})
	}

	return tags
}

type reflectSettings struct {
	index int
	tag   string
}
