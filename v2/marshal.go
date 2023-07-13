package extra

import (
	"reflect"
	"unsafe"

	"github.com/goccy/go-json"
	"github.com/gohort/extra/v2/internal/runtime"
)

// MarshalWithMapIndent is like MarshalWithMap but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalWithMapIndent(a any, m Map, prefix, indent string) ([]byte, error) {
	raw, err := json.MarshalIndentWithOption(a, prefix, indent, json.UnorderedMap())
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}

	return json.MarshalIndent(m, prefix, indent)
}

// MarshalWithMap is like Marshal but you can explicitly provide the [extra.Map]
// to the function.
func MarshalWithMap(a any, m Map) ([]byte, error) {
	return MarshalWithMapIndent(a, m, "", "")
}

// MarshalIndent is like Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(a any, prefix, indent string) ([]byte, error) {
	var m *Map

	if mapper, ok := (a).(Mapper); ok {
		m = mapper.GetExtraMap()
	} else {
		header := (*emptyInterface)(unsafe.Pointer(&a))
		typ := header.typ
		if typ.Kind() == reflect.Pointer {
			typ = typ.Elem()
		}

		// If this isn't a struct then just exit here with the json already marshalled.
		if typ.Kind() != reflect.Struct {
			return json.Marshal(a)
		}

		structHeader := (*runtime.StructType)(unsafe.Pointer(typ))
		for i := 0; i < len(structHeader.Fields); i++ {
			fieldHeader := structHeader.Fields[i]
			rtyp := (*runtime.Type)(unsafe.Pointer(fieldHeader.Typ))

			if runtime.PtrTo(rtyp).Implements(extraMapType) {
				m = (*Map)(add(header.ptr, fieldHeader.Offset))
				break
			}
		}
	}

	return MarshalWithMapIndent(a, *m, prefix, indent)
}

// Marshal returns the JSON encoding of v and flattens the Extra Map.
//
// Marshal traverses the value v recursively.
// If an encountered value implements the Marshaler interface
// and is not a nil pointer, Marshal calls its MarshalJSON method
// to produce JSON. If no MarshalJSON method is present but the
// value implements encoding.TextMarshaler instead, Marshal calls
// its MarshalText method and encodes the result as a JSON string.
// The nil pointer exception is not strictly necessary
// but mimics a similar, necessary exception in the behavior of
// UnmarshalJSON.
//
// UnmarshalWithMap uses goccy/go-json to Unmarshal JSON data.
//
// For more information please visit [goccy/go-json godoc page]
//
// [goccy/go-json godoc page]: https://pkg.go.dev/github.com/goccy/go-json#Unmarshal
func Marshal(a any) ([]byte, error) {
	return MarshalIndent(a, "", "")
}