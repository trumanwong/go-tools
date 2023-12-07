package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCache_WithValidOptions(t *testing.T) {
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	cache, err := NewCache(options, "test")

	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.Equal(t, "test", cache.prefix)
	assert.NotNil(t, cache.client)
}

func TestNewCache_WithInvalidOptions(t *testing.T) {
	options := &redis.Options{
		Addr:     "invalid_address",
		Password: "",
		DB:       0,
	}

	cache, err := NewCache(options, "test")

	assert.Error(t, err)
	assert.Nil(t, cache)
}
