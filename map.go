package extra

type Map interface {
	Set(key string, d interface{})
	Get(key string) interface{}
	GetOk(key string) (interface{}, bool)
	Delete(key string)
	Keys() []string
}
