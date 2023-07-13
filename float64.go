package extra

import "reflect"

var (
	_ Map = &Float64{}
)

type Float64 map[string]float64

func (a *Float64) Set(key string, d any) {
	switch val := d.(type) {
	case float32, float64:
		(*a)[key] = reflect.ValueOf(val).Float()
	}
}

func (a *Float64) Get(key string) any {
	return (*a)[key]
}

func (a *Float64) GetOk(key string) (any, bool) {
	d, ok := (*a)[key]
	return d, ok
}

func (a *Float64) Delete(key string) {
	delete(*a, key)
}

func (a *Float64) Keys() []string {
	keys := make([]string, len(*a))
	var i int
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}
