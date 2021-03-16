package store

// Store представляет собой key value хранилище.
// Сейчас реализована только in memory версия, так что ошибки можно игнорировать
type Store interface {
	Get(key string) (interface{}, bool, error)
	Set(key string, value interface{}) error
	Delete(key string) error
	List() (map[string]interface{}, error)
}
