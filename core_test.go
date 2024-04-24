package main

import (
	"errors"
	"testing"
)

func TestGet(t *testing.T) {
	const key = "key1"
	const value = "value1"

	var val interface{}
	var err error

	defer delete(store.m, value)

	val, err = Get(key)
	if err == nil {
		t.Error("expected an error")
	}
	if !errors.Is(err, ErrNoSuchKey) {
		t.Error("unexpected error: ", err)
	}
	store.m[key] = value

	val, err = Get(key)
	if err != nil {
		t.Error("expected no errors")
	}
	if val != value {
		t.Error("value mismatch", val, value)
	}
}

func TestDelete(t *testing.T) {
	const key = "key1"
	const value = "value1"

	var contains bool

	defer delete(store.m, value)

	store.m[key] = value
	_, contains = store.m[key]
	if !contains {
		t.Error("key/value not present")
	}
	Delete(key)

	_, contains = store.m[key]

	if contains {
		t.Error("delete failed")
	}
}

func TestPut(t *testing.T) {
	const key = "create-key"
	const value = "create-value"

	var val interface{}
	var contains bool

	defer delete(store.m, key)

	// Sanity check
	_, contains = store.m[key]
	if contains {
		t.Error("key/value already exists")
	}

	// err should be nil
	err := Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, contains = store.m[key]
	if !contains {
		t.Error("create failed")
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}
