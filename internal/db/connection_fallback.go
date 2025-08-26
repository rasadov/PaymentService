//go:build !js || !wasm

package db

import (
	"errors"
	"sync"
)

// Simple in-memory storage for non-WebAssembly builds
type kvStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func (k *kvStorage) Get(key string) (string, error) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	
	value, exists := k.data[key]
	if !exists {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (k *kvStorage) Put(key string, value string) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	if k.data == nil {
		k.data = make(map[string]string)
	}
	k.data[key] = value
	return nil
}

func (k *kvStorage) Delete(key string) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	
	delete(k.data, key)
	return nil
}
