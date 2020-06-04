package extra

var (
	_ Map = &Any{}
)

type Any map[string]interface{}

func (a *Any) Set(key string, d interface{}) {
	(*a)[key] = d
}

func (a *Any) Get(key string) interface{} {
	return (*a)[key]
}

func (a *Any) GetOk(key string) (interface{}, bool) {
	d, ok := (*a)[key]
	return d, ok
}

func (a *Any) Delete(key string) {
	delete(*a, key)
}

func (a *Any) Keys() []string {
	keys := make([]string, len(*a))
	var i int
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}
