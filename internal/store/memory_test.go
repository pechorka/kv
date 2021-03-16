package store

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemoryStore_GetSetDelete(t *testing.T) {
	store := NewMemory()

	const (
		testKey   = "testKey"
		testValue = "testValue"
	)

	val, ok, err := store.Get(testKey)
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, val)

	err = store.Set(testKey, testValue)
	require.NoError(t, err)

	val, ok, err = store.Get(testKey)
	require.NoError(t, err)
	require.True(t, ok)
	require.Exactly(t, testValue, val)

	err = store.Delete(testKey)
	require.NoError(t, err)

	val, ok, err = store.Get(testKey)
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, val)
}

func TestMemoryStore_List(t *testing.T) {
	store := NewMemory()

	testData := map[string]interface{}{
		"int":    1,
		"string": "string",
		"float":  1.0,
		"bool":   true,
	}

	for key, val := range testData {
		err := store.Set(key, val)
		require.NoError(t, err)
	}

	storeContent, err := store.List()
	require.NoError(t, err)

	for key, val := range testData {
		require.Exactly(t, val, storeContent[key])
	}

	// проверяем, что не можем менять содержимое хранилища через результат Memory.List()
	storeContent["pointerLickTest"] = "someVal"

	val, ok, err := store.Get("pointerLickTest")
	require.NoError(t, err)
	require.False(t, ok)
	require.Nil(t, val)
}
