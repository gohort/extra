package extra

type Map interface {
	Set(key string, d any)
	Get(key string) any
	GetOk(key string) (any, bool)
	Delete(key string)
	Keys() []string
}
