package extra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Unmarshal takes the pointer of an extra defined map and fills it with
// the fields that are not captured within the given `a` structure.
func Unmarshal(data []byte, a any, extras ...Map) error {
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
	if err := fillStruct(&tmp, a); err != nil {
		return err
	}

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

func fillStruct(extra Map, a any) error {
	tags := gatherTags(a)

	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return &json.InvalidUnmarshalError{Type: reflect.TypeOf(v)}
	}
	// Make a passthrough to handle embedded structures.
	for i := 0; i < v.Elem().NumField(); i++ {
		field := v.Elem().Field(i)
		fieldType := reflect.TypeOf(a).Elem()
		if field.Kind() == reflect.Struct {
			// If there's no tag found for this field then try unmarshalling it.
			if reflect.TypeOf(a).Elem().Field(i).Tag.Get("json") == "" {
				if field.IsValid() && field.CanSet() {
					if field.CanAddr() && field.CanInterface() {
						bb, err := json.Marshal(extra)
						if err != nil {
							return err
						}
						if err := json.Unmarshal(bb, field.Addr().Interface()); err != nil {
							return err
						}
						// Delete the field name if it's a capital casing or lower.
						fieldName := fieldType.Field(i).Name
						extra.Delete(fieldName)
						fieldName = fmt.Sprintf("%s%s", bytes.ToLower([]byte{fieldName[0]}), fieldName[1:])
						extra.Delete(fieldName)

						err = fillStruct(extra, field.Addr().Interface())
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}
	// Loop through tagged fields.
	for _, t := range tags {
		if val, ok := extra.GetOk(t.tag); ok {
			field := reflect.ValueOf(a).Elem().Field(t.index)
			if field.IsValid() && field.CanSet() {
				if field.CanAddr() && field.CanInterface() {
					bb, err := json.Marshal(val)
					if err != nil {
						return err
					}
					if err := json.Unmarshal(bb, field.Addr().Interface()); err != nil {
						return err
					}
				} else {
					fieldType := field.Type()
					valType := reflect.TypeOf(val)
					if !valType.AssignableTo(fieldType) || !valType.ConvertibleTo(fieldType) {
						return fmt.Errorf("%s can't go in %s: %w", valType, fieldType, ErrMismatchingTypes)
					}

					switch valType.Kind() {
					case reflect.Int:
						field.SetInt(val.(int64))
					case reflect.Bool:
						field.SetBool(val.(bool))
					case reflect.String:
						field.SetString(val.(string))
					case reflect.Float64:
						field.SetFloat(val.(float64))
					default:
						if field.CanAddr() && field.CanInterface() {
							bb, err := json.Marshal(val)
							if err != nil {
								return err
							}
							json.Unmarshal(bb, field.Addr().Interface())
						}
					}
				}

				extra.Delete(t.tag)
			}
		}
	}

	return nil
}
