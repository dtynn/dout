package storage

type Storage interface {
	Set(key, val interface{}) error
	Setex(key, val interface{}, expire int64) error
	Get(key interface{}) (interface{}, error)
	Del(key interface{}) error
}
