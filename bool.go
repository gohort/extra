package extra

var (
	_ Map = &Bools{}
)

type Bools map[string]bool

func (a *Bools) Set(key string, d interface{}) {
	if b, ok := d.(bool); ok {
		(*a)[key] = b
	}
}

func (a *Bools) Get(key string) interface{} {
	return (*a)[key]
}

func (a *Bools) GetOk(key string) (interface{}, bool) {
	d, ok := (*a)[key]
	return d, ok
}

func (a *Bools) Delete(key string) {
	delete(*a, key)
}

func (a *Bools) Keys() []string {
	keys := make([]string, len(*a))
	var i int
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}
