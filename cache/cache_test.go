package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewCache_WithValidOptions(t *testing.T) {
	options := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
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

func TestCache_LRemBeforeKey(t *testing.T) {
	options := &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	cache, err := NewCache(options, "")

	err = cache.LRemBeforeKey(context.Background(), &LRemByValueRequest{
		Key:    "test",
		Value:  "100",
		Prefix: nil,
	})
	assert.NoError(t, err)
}
