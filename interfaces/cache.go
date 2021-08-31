package interfaces

type CacheInterface interface {
	Save(key string, value interface{})
	Get(key string) interface{}

}

