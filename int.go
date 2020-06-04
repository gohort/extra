package extra

import "reflect"

var (
	_ Map = &Ints{}
)

type Ints map[string]int

func (a *Ints) Set(key string, d interface{}) {
	switch val := d.(type) {
	case uint, uint8, uint16, uint32,
		uint64, int, int8, int16, int32, int64:
		(*a)[key] = int(reflect.ValueOf(val).Int())
	case float32, float64:
		(*a)[key] = int(reflect.ValueOf(val).Float())
	}
}

func (a *Ints) Get(key string) interface{} {
	return (*a)[key]
}

func (a *Ints) GetOk(key string) (interface{}, bool) {
	d, ok := (*a)[key]
	return d, ok
}

func (a *Ints) Delete(key string) {
	delete(*a, key)
}

func (a *Ints) Keys() []string {
	keys := make([]string, len(*a))
	var i int
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}
