package main

import (
	"errors"
	"sync"
)

type LockableMap struct {
	sync.RWMutex
	m map[string]string
}

var store = LockableMap{
	m: make(map[string]string),
}

// store        = make(map[string]string)
var ErrNoSuchKey = errors.New("no such key")

func Put(key, value string) error {
	store.Lock()
	defer store.Unlock()

	store.m[key] = value
	return nil
}

func Get(key string) (string, error) {
	store.Lock()
	defer store.Unlock()

	value, ok := store.m[key]
	if !ok {
		return "", ErrNoSuchKey
	}
	return value, nil
}

func Delete(key string) error {
	store.Lock()
	defer store.Unlock()

	delete(store.m, key)
	return nil
}

func HealthCheck() bool {
	return true
}
