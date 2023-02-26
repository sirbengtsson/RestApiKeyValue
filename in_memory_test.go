package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSettingValueShouldBeOk(t *testing.T) {
	keyValueStore := &MemoryStore{}

	value := "testValue"
	key := "testKey"
	keyValueStore.setValue(key, value)
	var result, err = keyValueStore.readValue(key)

	assert.Equal(t, result, value)
	assert.Equal(t, err, nil)
}

func TestDeletionValueShouldBeOk(t *testing.T) {
	keyValueStore := &MemoryStore{}

	value := "testValue"
	key := "testKey"

	keyValueStore.setValue(key, value)

	keyValueStore.deleteValue(key)
	var result, err = keyValueStore.readValue(key)

	assert.Equal(t, result, "")
	assert.Errorf(t, err, "key not found: %s", key)
}

func TestDoubleDeletionShouldFail(t *testing.T) {
	keyValueStore := &MemoryStore{}

	value := "testValue"
	key := "testKey"

	keyValueStore.setValue(key, value)

	keyValueStore.deleteValue(key)
	doubleDeleteErr := keyValueStore.deleteValue(key)

	assert.Errorf(t, doubleDeleteErr, "key not found: %s", key)
}
