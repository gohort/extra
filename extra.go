package extra

import (
	"encoding/json"
	"reflect"
)

// Unmarshal takes the pointer of an extra defined map and fills it with
// the fields that are not captured within the given `a` structure.
func Unmarshal(data []byte, a interface{}, extras ...Map) error {
	switch {
	case a == nil:
		return ErrNilInterface
	}
	// Create a temporary map to capture all the elements of the JSON.
	var tmp = make(Any)
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	// Fill the structure provided with whatever fields exist in the temp map.
	fillStruct(&tmp, a)
	for _, ex := range extras {
		// If the map given is nil, then make the map.
		val := reflect.ValueOf(ex).Elem()
		if val.IsNil() {
			val.Set(reflect.MakeMap(val.Type()))
		}
		// Set what can be set into the given extra map based on its type.
		for k, v := range tmp {
			ex.Set(k, v)
		}
	}

	return nil
}

func fillStruct(extra Map, a interface{}) {
	tags := gatherTags(a)

	for _, t := range tags {
		if val, ok := extra.GetOk(t.tag); ok {
			field := reflect.ValueOf(a).Elem().Field(t.index)
			if field.IsValid() && field.CanSet() {
				switch reflect.TypeOf(val).Kind() {
				case reflect.Int:
					field.SetInt(val.(int64))
				case reflect.Bool:
					field.SetBool(val.(bool))
				case reflect.String:
					field.SetString(val.(string))
				case reflect.Float64:
					field.SetFloat(val.(float64))
				}

				extra.Delete(t.tag)
			}
		}
	}
}

// Marshal takes all the fields of `a` and inserts them into the extra map.
// Then marshals the map.
func Marshal(a interface{}, extras ...Map) ([]byte, error) {
	switch {
	case a == nil:
		return nil, ErrNilInterface
	}

	var (
		// Get all the tags in the structure provided.
		tags = gatherTags(a)
		// Create a copy map that will contain all the elements in extras.
		extra = make(Any)
	)

	for _, ex := range extras {
		// Get the keys from the map.
		keys := ex.Keys()
		// Loop through all the keys and copy the values within the map.
		// This is to ensure that we do not change the actual values within
		// the map that is given to us.
		for _, k := range keys {
			extra[k] = ex.Get(k)
		}
		// Loop through all the tags in the structure provided.
		for _, t := range tags {
			field := reflect.ValueOf(a).Elem().Field(t.index)
			if field.IsValid() {
				switch field.Kind() {
				case reflect.Int:
					extra.Set(t.tag, field.Int())
				case reflect.String:
					extra.Set(t.tag, field.String())
				case reflect.Float64:
					extra.Set(t.tag, field.Float())
				case reflect.Bool:
					extra.Set(t.tag, field.Bool())
				default:
					if field.CanInterface() {
						extra.Set(t.tag, field.Interface())
					}
				}
			}
		}
	}

	return json.Marshal(extra)
}
