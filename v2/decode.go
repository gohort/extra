package extra

import (
	"fmt"
	"reflect"
)

//
// This file is still under construction.
//
// Thinking of making something that's akin to mapstructure's two step
// unmarshal/decode if you're used to doing it.
//

// Decode
//
// Decode assumes that [Map] has values and is non-nil.
//
// > Example Usage:
//
//	type User struct {
//		Username string `json:"username"`
//		X extra.Map `json:"-"`
//	}
//	var u User
//	if err := json.Unmarshal(data, &u.X); err != nil {
//		// handle err
//	}
//
//	if err := extra.Decode(&u.X, &u); err != nil {
//		// handle err
//	}
// func Decode(m *Map, a any) error {
// 	header := (*emptyInterface)(unsafe.Pointer(&a))
// 	if err := validateType(header.typ, uintptr(header.ptr)); err != nil {
// 		return err
// 	}
// 	if m == nil {
// 		return ErrNilMap
// 	}

// 	// reflect.ValueOf(a).Elem()
// 	// elem := header.typ.Elem()
// 	elem := header.typ
// 	if elem.Elem().Kind() != reflect.Invalid {
// 		// fmt.Printf("kin[%s]\n", elem.Elem().Kind().String())
// 		// runtime.RType2Type(elem.Elem()).

// 		elem = elem.Elem()
// 		// if elem.Kind() == reflect.Interface {
// 		// 	elem =
// 		// }
// 	}

// 	for i := 0; i < elem.NumField(); i++ {
// 		if !elem.Field(i).IsExported() {
// 			continue
// 		}

// 		if name, ignore := parseTagName(elem.Field(i)); !ignore {
// 			if val, ok := (*m)[name]; ok {
// 				// Set in struct
// 				field := reflect.NewAt(elem.Field(i).Type, add(header.ptr, elem.Field(i).Offset))
// 				if field.Kind() == reflect.Pointer {
// 					field = field.Elem()
// 				}

// 				// unwrappedInterface, unwrappedKind := unwrapKind(field)
// 				// switch unwrappedKind {
// 				// case reflect.Struct:
// 				// 	Decode(m, &unwrappedInterface)
// 				// 	continue
// 				// }
// 				// fmt.Printf("kin[%s]\n", field.Kind().String())
// 				if !elem.Field(i).Type.AssignableTo(reflect.TypeOf(val)) {
// 					fmt.Printf("%s not assignable to %s\n", elem.Field(i).Type.String(), reflect.TypeOf(val))
// 					// fmt.Printf("kin[%s]\n", field.Elem().Kind().String())
// 					// fmt.Printf("*knd=%s\n", field.Kind())
// 					// field.Type().Kind()
// 					// fmt.Printf("*typ=%s\n", field.Type().String())
// 					// Decode(m, (any)(add(header.ptr, elem.Field(i).Offset)))
// 					// v := field.Interface()
// 					// field.Elem().SetZero()
// 					dataVal := reflect.Indirect(field.Elem())
// 					if err := decodeStruct(m, dataVal); err != nil {
// 						return err
// 					}
// 					// v := reflect.ValueOf(a).Elem().Field(i)
// 					// Decode(m, (any)(&v))
// 					continue
// 				} else {
// 					fmt.Printf(" typ=%s\n", field.Type().Name())
// 				}

// 				field.Set(reflect.ValueOf(val))
// 				// *(*any)(add(header.ptr, elem.Field(i).Offset)) = val

// 				// Remove from map
// 				delete(*m, name)
// 			}
// 		}
// 		// add(header.ptr, header.typ.Elem().Field(i).Offset)
// 	}

// 	return nil
// }

//lint:ignore U1000 Ignore unused function temporarily for debugging
func decodeStruct(m *Map, v reflect.Value) error {
	// if v.Kind() != reflect.Struct {
	// 	// panic("invalid kind: not struct")
	// 	return fmt.Errorf("%s: %w", v.Kind(), errors.New("not struct"))
	// }

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Pointer {
			if err := decodeStruct(m, field); err != nil {
				// panic("recursive call")
				return fmt.Errorf("recursive decodeStruct: %w", err)
			}
			continue
		}


	}

	return nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func unwrapKind(v reflect.Value) (any, reflect.Kind) {
	kind := v.Kind()
	switch kind {
	case reflect.Pointer:
		return unwrapKind(v.Elem())
	case reflect.Map, reflect.Slice, reflect.Struct:
		return (any)(v.UnsafePointer()), kind
	}
	return v.Interface(), v.Kind()
}