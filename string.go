package extra

var (
	_ Map = &Strings{}
)

type Strings map[string]string

func (a *Strings) Set(key string, d any) {
	if str, ok := d.(string); ok {
		(*a)[key] = str
	}
}

func (a *Strings) Get(key string) any {
	return (*a)[key]
}

func (a *Strings) GetOk(key string) (any, bool) {
	d, ok := (*a)[key]
	return d, ok
}

func (a *Strings) Delete(key string) {
	delete(*a, key)
}

func (a *Strings) Keys() []string {
	keys := make([]string, len(*a))
	var i int
	for k := range *a {
		keys[i] = k
		i++
	}
	return keys
}
