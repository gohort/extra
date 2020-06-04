package extra

import "reflect"

var (
	_ Map = &Float32{}
)

type Float32 map[string]float32

func (a *Float32) Set(key string, d interface{}) {
	switch val := d.(type) {
	case float32, float64:
		(*a)[key] = float32(reflect.ValueOf(val).Float())
	}
}

func (a *Float32) Get(key string) interface{} {
	return (*a)[key]
}

func (a *Float32) GetOk(key string) (interface{}, bool) {
	d, ok := (*a)[key]
	return d, ok
}

func (a *Float32) Delete(key string) {
	delete(*a, key)
}

func (a *Float32) Keys() []string {
	keys := make([]string, len(*a))
	var i int
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}
