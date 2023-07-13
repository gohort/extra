package extra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Marshal takes all the fields of `a` and inserts them into the extra map.
// Then marshals the map.
func Marshal(a any, extras ...Map) ([]byte, error) {
	switch {
	case a == nil:
		return nil, ErrNilInterface
	}

	var (
		extra = make(Any)
	)

	for _, ex := range extras {
		keys := ex.Keys()
		for _, key := range keys {
			extra.Set(key, ex.Get(key))
		}
	}

	if err := MarshalInto(a, extra); err != nil {
		return nil, err
	}

	return json.Marshal(extra)
}

func MarshalInto(a any, extra Any) error {
	v := reflect.ValueOf(a).Elem()
	t := reflect.TypeOf(a).Elem()

	for i := 0; i < v.NumField(); i++ {
		var key string
		key = t.Field(i).Tag.Get("json")
		if key == "-" {
			continue
		}
		if key == "" {
			name := t.Field(i).Name
			key = fmt.Sprintf("%s%s", bytes.ToLower([]byte{name[0]}), name[1:])
		}

		field := v.Field(i)
		if field.IsValid() {
			switch field.Kind() {
			case reflect.Int:
				extra.Set(key, field.Int())
			case reflect.String:
				extra.Set(key, field.String())
			case reflect.Float64:
				extra.Set(key, field.Float())
			case reflect.Bool:
				extra.Set(key, field.Bool())
			default:
				if field.CanInterface() {
					// If the type is a structure, then recurse this function.
					if field.CanAddr() && field.Kind() == reflect.Struct {
						if err := MarshalInto(field.Addr().Interface(), extra); err != nil {
							return err
						}
					} else {
						if field.Kind() == reflect.Slice { // If it's a slice, then append to map.
							extra.Set(key, field.Interface())
						} else { // Maps should land here
							bb, err := json.Marshal(field.Interface())
							if err != nil {
								return err
							}
							if err := json.Unmarshal(bb, &extra); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}
