package store

import "sync"

type Memory struct {
	sync.RWMutex
	store map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{
		store: make(map[string]interface{}),
	}
}

// Get возвращает значение для ключа
func (kv *Memory) Get(key string) (interface{}, bool, error) {
	kv.RLock()
	defer kv.RUnlock()

	val, ok := kv.store[key]
	return val, ok, nil
}

// Set сохраняет значение val под ключом key
func (kv *Memory) Set(key string, val interface{}) error {
	kv.Lock()
	defer kv.Unlock()

	kv.store[key] = val

	return nil
}

// Delete удаляет значение под ключом key
func (kv *Memory) Delete(key string) error {
	kv.Lock()
	defer kv.Unlock()

	delete(kv.store, key)

	return nil
}

// List возвращает содержимое хранилища
func (kv *Memory) List() (map[string]interface{}, error) {
	kv.RLock()
	defer kv.RUnlock()

	// копируем map. Если вернуть просто kv.store, то мы выдадим указатель на map.
	// Тогда у нас не будет гарантий корректности данных
	storeContent := make(map[string]interface{})
	for key, val := range kv.store {
		storeContent[key] = val
	}

	return storeContent, nil
}
