package main

import (
	"fmt"
	"sync"
)

type MemoryStore struct {
}

var cache = make(map[string]string)
var mutex = &sync.RWMutex{}

func (store *MemoryStore) setValue(key string, value string) error {
	mutex.Lock()
	cache[key] = value
	mutex.Unlock()

	return nil
}

func (store *MemoryStore) readValue(key string) (string, error) {
	mutex.RLock()
	value, ok := cache[key]
	mutex.RUnlock()

	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

func (store *MemoryStore) deleteValue(key string) error {
	mutex.Lock()

	_, ok := cache[key]

	if ok {
		delete(cache, key)
	} else {
		return fmt.Errorf("key not found: %s", key)
	}

	mutex.Unlock()

	return nil
}

func (store *MemoryStore) getAllKeys() ([]string, error) {
	mutex.RLock()
	keys := make([]string, 0, len(cache))

	for key := range cache {
		keys = append(keys, key)
	}
	mutex.RUnlock()

	return keys, nil
}
