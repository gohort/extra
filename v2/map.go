package extra

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/gohort/extra/v2/internal/runtime"
)

var _ Mapper = &Map{}

type Map map[string]any

func (m *Map) GetExtraMap() *Map {
	return m
}

func (m *Map) GetExtraField(k string) (any, bool) {
	v, ok := (*m)[k]
	return v, ok
}

func (m *Map) RemoveExtraKey(key string) {
	delete(*m, key)
}

// emptyInterface is implemented in the reflect package and looks very similar
// to this one, the difference being is that we can utilize it here to skip
// any extra allocations that the reflect package does when taking values of
// passed variables and instead inspect the type ourselves.
//
// Credit: https://github.com/goccy/go-json/blob/master/decode.go#L23
type emptyInterface struct {
	typ *runtime.Type
	ptr unsafe.Pointer
}

// check that this type and pointer are indeed a pointer and non-nil.
func validateType(typ *runtime.Type, p uintptr) error {
	if typ == nil || typ.Kind() != reflect.Ptr || p == 0 {
		return &ErrInvalidUnmarshal{Type: runtime.RType2Type(typ)}
	}
	return nil
}

func parseTagName(f reflect.StructField) (name string, ignore bool) {
	tag := f.Tag.Get("json")

	if tag == "" {
		return f.Name, false
	}

	if tag == "-" {
		return "", true
	}

	if i := strings.Index(tag, ","); i != -1 {
		if i == 0 {
			return f.Name, false
		}
		return tag[:i], false
	}

	return tag, false
}

func parseTagNameRField(f *runtime.StructField) (name string, ignore bool) {
	// tag := f.Tag.Get("json")
	tag := f.Tag().Get("json")

	if tag == "" {
		return f.Name.Name(), false
	}

	if tag == "-" {
		return "", true
	}

	if i := strings.Index(tag, ","); i != -1 {
		if i == 0 {
			return f.Name.Name(), false
		}
		return tag[:i], false
	}

	return tag, false
}
