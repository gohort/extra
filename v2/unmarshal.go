package extra

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/goccy/go-json"
	"github.com/gohort/extra/v2/internal/runtime"
)

var (
	extraMapType = reflect.TypeOf((*Mapper)(nil)).Elem()
)

// UnmarshalWithMap parses the JSON-encoded data and stores the result
// in the value pointed to by v and m. If v is nil or not a pointer,
// UnmarshalWithMap returns an ErrInvalidUnmarshal. If m is nil,
// UnmarshalWithMap returns an ErrNilMap.
//
// UnmarshalWithMap uses the inverse of the encodings that
// Marshal uses, allocating maps, slices, and pointers as necessary,
// with the following additional rules:
//
// UnmarshalWithMap uses goccy/go-json to Unmarshal JSON data.
//
// For more information please visit [goccy/go-json godoc page]
//
// [goccy/go-json godoc page]: https://pkg.go.dev/github.com/goccy/go-json#Unmarshal
func UnmarshalWithMap(data []byte, v any, m *Map) error {
	// Extract the header meta data of this variable to inspect its type.
	header := (*emptyInterface)(unsafe.Pointer(&v))
	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
		return err
	}
	if m == nil {
		return ErrNilMap
	}

	// If this is a reference to a pointer then deref and unmarshal.
	switch header.typ.Elem().Kind() {
	case reflect.Pointer:
		return UnmarshalWithMap(data, reflect.ValueOf(v).Elem().Interface(), m)
	case reflect.Struct:
	default:
		return json.Unmarshal(data, v)
	}

	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("unmarshal map: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal obj: %w", err)
	}

	elem := header.typ.Elem()
	for i := 0; i < elem.NumField(); i++ {
		if tag, ignore := parseTagName(elem.Field(i)); !ignore {
			delete(*m, tag)
		}
	}

	return nil
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v and m. If v is nil or not a pointer,
// Unmarshal returns an ErrInvalidUnmarshal. If m is nil,
// Unmarshal returns an ErrNilMap.
//
// Unmarshal uses the inverse of the encodings that
// Marshal uses, allocating maps, slices, and pointers as necessary,
// with the following additional rules:
//
// Unmarshal uses goccy/go-json to Unmarshal JSON data.
//
// For more information please visit [goccy/go-json godoc page]
//
// [goccy/go-json godoc page]: https://pkg.go.dev/github.com/goccy/go-json#Unmarshal
func Unmarshal(data []byte, v any) error {
	// Extract the header meta data of this variable to inspect its type.
	header := (*emptyInterface)(unsafe.Pointer(&v))
	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
		return err
	}

	// If this is a reference to a pointer then deref and unmarshal.
	switch header.typ.Elem().Kind() {
	case reflect.Pointer:
		return Unmarshal(data, reflect.ValueOf(v).Elem().Interface())
	case reflect.Struct:
	default:
		return json.Unmarshal(data, v)
	}

	// Unmarshal the data into the object given.
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	// Either get the mapper from an embedded Map object or make a new one.
	var m *Map
	var mapSet bool
	if mapper, ok := v.(Mapper); ok {
		// Get a reference to the object's Map.
		m = mapper.GetExtraMap()
		mapSet = true
		*m = make(Map) // initialize the object's map.
	} else {
		// Since we don't know where the object's Map is right away we have to
		// initialize a new one then set it as a reference later.
		newMap := make(Map)
		m = &newMap
	}

	// Unmarshal the data into our extra Map.
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	// Dissect the struct header from our header type.
	structHeader := (*runtime.StructType)(unsafe.Pointer(header.typ.Elem()))
	// Loop through our fields.
	for i := 0; i < len(structHeader.Fields); i++ {
		fieldHeader := structHeader.Fields[i]
		rtyp := (*runtime.Type)(unsafe.Pointer(fieldHeader.Typ))

		// If this field isn't exported then don't bother parsing its tags.
		if !fieldHeader.Name.IsExported() {
			continue
		}

		if name, ignore := parseTagNameRField(&fieldHeader); !ignore {
			delete(*m, name)
			continue
		}

		// If we already have a reference to this object's Map then don't continue
		// any further.
		if mapSet {
			continue
		}

		// Check if this field implements the extraMap interface.
		if runtime.PtrTo(rtyp).Implements(extraMapType) {
			*(*Map)(add(header.ptr, fieldHeader.Offset)) = *m
		}
	}

	return nil
}

// add returns p+x.
//
// This function was taken from the [reflect package]
//
// [reflect package]: https://cs.opensource.google/go/go/+/refs/tags/go1.20.5:src/reflect/type.go;l=1100
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}
