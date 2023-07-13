package runtime

import (
	"strconv"
	"unsafe"
)

// Rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// Rtype must be kept in sync with ../runtime/type.go:/^type._type.
type Rtype struct {
	size       uintptr
	Ptrdata    uintptr // number of bytes in the type that can contain pointers
	// hash       uint32  // hash of type; avoids computation in hash tables
	// tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	// kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal     func(unsafe.Pointer, unsafe.Pointer) bool
	Gcdata    *byte   // garbage collection data
	Str       nameOff // string form
	PtrToThis typeOff // type for pointer to this type, may be zero
}

// Struct field
type StructField struct {
	Name   name    // name is always non-empty
	Typ    *Rtype  // type of field
	Offset uintptr // byte offset of field
}

func (s *StructField) Tag() StructTag {
	return StructTag(s.Name.tag())
}

// func (f *StructField) embedded() bool {
// 	return f.Name.embedded()
// }

// StructType represents a struct type.
// type StructType struct {
// 	Type
// 	pkgPath name
// 	fields  []StructField // sorted by offset
// }
type StructType struct {
	Rtype
	pkgPath name
	Fields  []StructField // sorted by offset
}

// name is an encoded type name with optional extra data.
//
// The first byte is a bit field containing:
//
//	1<<0 the name is exported
//	1<<1 tag data follows the name
//	1<<2 pkgPath nameOff follows the name and tag
//	1<<3 the name is of an embedded (a.k.a. anonymous) field
//
// Following that, there is a varint-encoded length of the name,
// followed by the name itself.
//
// If tag data is present, it also has a varint-encoded length
// followed by the tag itself.
//
// If the import path follows, then 4 bytes at the end of
// the data form a nameOff. The import path is only set for concrete
// methods that are defined in a different package than their type.
//
// If a name starts with "*", then the exported bit represents
// whether the pointed to type is exported.
//
// Note: this encoding must match here and in:
//   cmd/compile/internal/reflectdata/reflect.go
//   runtime/type.go
//   internal/reflectlite/type.go
//   cmd/link/internal/ld/decodesym.go
type name struct {
	bytes *byte
}

func (n name) data(off int, whySafe string) *byte {
	return (*byte)(add(unsafe.Pointer(n.bytes), uintptr(off), whySafe))
}

func (n name) IsExported() bool {
	return (*n.bytes)&(1<<0) != 0
}

func (n name) hasTag() bool {
	return (*n.bytes)&(1<<1) != 0
}

func (n name) embedded() bool {
	return (*n.bytes)&(1<<3) != 0
}

// readVarint parses a varint as encoded by encoding/binary.
// It returns the number of encoded bytes and the encoded value.
func (n name) readVarint(off int) (int, int) {
	v := 0
	for i := 0; ; i++ {
		x := *n.data(off+i, "read varint")
		v += int(x&0x7f) << (7 * i)
		if x&0x80 == 0 {
			return i + 1, v
		}
	}
}

// writeVarint writes n to buf in varint form. Returns the
// number of bytes written. n must be nonnegative.
// Writes at most 10 bytes.
func writeVarint(buf []byte, n int) int {
	for i := 0; ; i++ {
		b := byte(n & 0x7f)
		n >>= 7
		if n == 0 {
			buf[i] = b
			return i + 1
		}
		buf[i] = b | 0x80
	}
}

type UnsafeString struct {
	Data unsafe.Pointer
	Len  int
}

func (n name) Name() (s string) {
	if n.bytes == nil {
		return ""
	}
	i, l := n.readVarint(1)
	hdr := (*UnsafeString)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(n.data(1+i, "non-empty string"))
	hdr.Len = l
	return
}

func (n name) tag() (s string) {
	if !n.hasTag() {
		return ""
	}
	i, l := n.readVarint(1)
	i2, l2 := n.readVarint(1 + i + l)
	hdr := (*UnsafeString)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(n.data(1+i+l+i2, "non-empty string"))
	hdr.Len = l2
	return
}

type nameOff int32 // offset to a name
type typeOff int32 // offset to an *Rtype

func (t *Rtype) Size() uintptr { return t.size }

func (t *Rtype) Align() int { return int(t.align) }

func (t *Rtype) FieldAlign() int { return int(t.fieldAlign) }

// add returns p+x.
//
// The whySafe string is ignored, so that the function still inlines
// as efficiently as p+x, but all call sites should use the string to
// record why the addition is safe, which is to say why the addition
// does not cause x to advance to the very end of p's allocation
// and therefore point incorrectly at the next block in memory.
func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type StructTag string

// Get returns the value associated with key in the tag string.
// If there is no such key in the tag, Get returns the empty string.
// If the tag does not have the conventional format, the value
// returned by Get is unspecified. To determine whether a tag is
// explicitly set to the empty string, use Lookup.
func (tag StructTag) Get(key string) string {
	v, _ := tag.Lookup(key)
	return v
}

// Lookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string. If the tag does not have the conventional format,
// the value returned by Lookup is unspecified.
func (tag StructTag) Lookup(key string) (value string, ok bool) {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value, true
		}
	}
	return "", false
}
