package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

// Cache is a struct that represents a Redis cache.
// It contains a prefix that is prepended to all keys in the cache,
// and a client that is used to interact with the Redis server.
type Cache struct {
	prefix string
	client *redis.Client
}

// NewCache is a function that creates a new Cache.
// It takes a pointer to a redis.Options struct, which contains options for the Redis client,
// and a prefix string, which is prepended to all keys in the cache.
// It returns a pointer to the created Cache and an error.
func NewCache(options *redis.Options, prefix string) (*Cache, error) {
	return &Cache{
		prefix: prefix,
		client: redis.NewClient(options),
	}, nil
}

// prefixKey is a method of Cache that prepends the cache's prefix to a key.
// If a custom prefix is provided, it is used instead of the cache's prefix.
func (c *Cache) prefixKey(key string, prefix *string) string {
	if prefix != nil {
		return *prefix + key
	}
	return c.prefix + key
}

// SetCacheRequest is a struct that represents a request to set a value in the cache.
// It contains the key to set, the value to set it to, the expiration time in seconds,
// and an optional custom prefix.
type SetCacheRequest struct {
	Key     string
	Value   []byte
	Seconds int64
	Prefix  *string
}

// Set is a method of Cache that sets a value in the cache.
// It takes a context and a pointer to a SetCacheRequest struct,
// and returns an error.
// The method uses the Redis SET command to set the value.
func (c *Cache) Set(ctx context.Context, request *SetCacheRequest) error {
	_, err := c.client.Set(ctx, c.prefixKey(request.Key, request.Prefix), request.Value, time.Duration(request.Seconds)*time.Second).Result()
	return err
}

// GetCacheRequest is a struct that represents a request to get a value from the cache.
// It contains the key to get and an optional custom prefix.
type GetCacheRequest struct {
	Key    string
	Prefix *string
}

// Get is a method of Cache that gets a value from the cache.
// It takes a context and a pointer to a GetCacheRequest struct,
// and returns the value as a string and an error.
// The method uses the Redis GET command to get the value.
func (c *Cache) Get(ctx context.Context, request *GetCacheRequest) (string, error) {
	return c.client.Get(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// IncRequest is a struct that represents a request to increment a value in the cache.
// It contains the key to increment and an optional custom prefix.
type IncRequest struct {
	Key    string
	Prefix *string
}

// Inc is a method of Cache that increments a value in the cache.
// It takes a context and a pointer to an IncRequest struct,
// and returns the new value as an int64 and an error.
// The method uses the Redis INCR command to increment the value.
func (c *Cache) Inc(ctx context.Context, request *IncRequest) (int64, error) {
	return c.client.Incr(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// IncrByRequest is a struct that represents a request to increment a value in the cache by a certain amount.
// It contains the key to increment, the amount to increment by, and an optional custom prefix.
type IncrByRequest struct {
	Key    string
	Value  int64
	Prefix *string
}

// IncrBy is a method of Cache that increments a value in the cache by a certain amount.
// It takes a context and a pointer to an IncrByRequest struct,
// and returns the new value as an int64 and an error.
// The method uses the Redis INCRBY command to increment the value.
func (c *Cache) IncrBy(ctx context.Context, request *IncrByRequest) (int64, error) {
	return c.client.IncrBy(ctx, c.prefixKey(request.Key, request.Prefix), request.Value).Result()
}

// DeleteRequest is a struct that represents a request to delete a key from the cache.
// It contains the key to delete and an optional custom prefix.
type DeleteRequest struct {
	Key    string
	Prefix *string
}

// Delete is a method of Cache that deletes a key from the cache.
// It takes a context and a pointer to a DeleteRequest struct,
// and returns the number of keys that were deleted as an int64 and an error.
// The method uses the Redis DEL command to delete the key.
func (c *Cache) Delete(ctx context.Context, request *DeleteRequest) (int64, error) {
	return c.client.Del(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// LPushRequest is a struct that represents a request to push a value onto a list in the cache.
// It contains the key of the list, the value to push, and an optional custom prefix.
type LPushRequest struct {
	Key    string
	Value  []interface{}
	Prefix *string
}

// LPush is a method of Cache that pushes a value onto a list in the cache.
// It takes a context and a pointer to a LPushRequest struct,
// and returns the length of the list after the push as an int64 and an error.
// The method uses the Redis LPUSH command to push the value.
func (c *Cache) LPush(ctx context.Context, request *LPushRequest) (int64, error) {
	return c.client.LPush(ctx, c.prefixKey(request.Key, request.Prefix), request.Value...).Result()
}

// LPopRequest is a struct that represents a request to pop a value from a list in the cache.
// It contains the key of the list and an optional custom prefix.
type LPopRequest struct {
	Key    string
	Prefix *string
}

// LPop is a method of Cache that pops a value from a list in the cache.
// It takes a context and a pointer to a LPopRequest struct,
// and returns the popped value as a string and an error.
// The method uses the Redis LPOP command to pop the value.
func (c *Cache) LPop(ctx context.Context, request *LPopRequest) (string, error) {
	return c.client.LPop(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// RPopRequest is a struct that represents a request to pop a value from the end of a list in the cache.
// It contains the key of the list and an optional custom prefix.
type RPopRequest struct {
	Key    string
	Prefix *string
}

// RPop is a method of Cache that pops a value from the end of a list in the cache.
// It takes a context and a pointer to a RPopRequest struct,
// and returns the popped value as a string and an error.
// The method uses the Redis RPOP command to pop the value.
func (c *Cache) RPop(ctx context.Context, request *RPopRequest) (string, error) {
	return c.client.RPop(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// LRangeRequest is a struct that represents a request to get a range of values from a list in the cache.
// It contains the key of the list, the range to get, and an optional custom prefix.
type LRangeRequest struct {
	Key    string
	Value  string
	Prefix *string
}

// LRange is a method of Cache that gets a range of values from a list in the cache.
// It takes a context and a pointer to a LRangeRequest struct,
// and returns the values as a slice of strings and an error.
// The method uses the Redis LRANGE command to get the values.
func (c *Cache) LRange(ctx context.Context, request *LRangeRequest) ([]string, error) {
	return c.client.LRange(ctx, c.prefixKey(request.Key, request.Prefix), 0, -1).Result()
}

// GetListIndexRequest is a struct that represents a request to get the index of a value in a list in the cache.
// It contains the key of the list, the value to find, and an optional custom prefix.
type GetListIndexRequest struct {
	Key    string
	Value  string
	Prefix *string
}

// GetListIndex is a method of Cache that gets the index of a value in a list in the cache.
// It takes a context and a pointer to a GetListIndexRequest struct,
// and returns the index of the value as an int and an error.
// The method uses the Redis LRANGE command to get the list, and then iterates over the list to find the value.
func (c *Cache) GetListIndex(ctx context.Context, request *GetListIndexRequest) (int, error) {
	list, err := c.client.LRange(ctx, c.prefixKey(request.Key, request.Prefix), 0, -1).Result()
	if err != nil {
		return 0, err
	}
	for i, v := range list {
		if v == request.Value {
			return len(list) - i, nil
		}
	}
	return 0, errors.New("no found")
}

// TTLRequest is a struct that represents a request to get the time to live of a key in the cache.
// It contains the key and an optional custom prefix.
type TTLRequest struct {
	Key    string
	Prefix *string
}

// TTL is a method of Cache that gets the time to live of a key in the cache.
// It takes a context and a pointer to a TTLRequest struct,
// and returns the time to live as a time.Duration and an error.
// The method uses the Redis TTL command to get the time to live.
func (c *Cache) TTL(ctx context.Context, request *TTLRequest) (time.Duration, error) {
	return c.client.TTL(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// ExistsRequest is a struct that represents a request to check if a key exists in the cache.
// It contains the key and an optional custom prefix.
type ExistsRequest struct {
	Key    string
	Prefix *string
}

// Exists is a method of Cache that checks if a key exists in the cache.
// It takes a context and a pointer to an ExistsRequest struct,
// and returns a boolean indicating whether the key exists and an error.
// The method uses the Redis EXISTS command to check if the key exists.
func (c *Cache) Exists(ctx context.Context, request *ExistsRequest) (bool, error) {
	val, err := c.client.Exists(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

// ExpireRequest is a struct that represents a request to set the expiration time of a key in the cache.
// It contains the key, the expiration time in seconds, and an optional custom prefix.
type ExpireRequest struct {
	Key     string
	Seconds int64
	Prefix  *string
}

// Expire is a method of Cache that sets the expiration time of a key in the cache.
// It takes a context and a pointer to an ExpireRequest struct,
// and returns a boolean indicating whether the expiration time was set and an error.
// The method uses the Redis EXPIRE command to set the expiration time.
func (c *Cache) Expire(ctx context.Context, request *ExpireRequest) (bool, error) {
	return c.client.Expire(ctx, c.prefixKey(request.Key, request.Prefix), time.Duration(request.Seconds)*time.Second).Result()
}

// ZRangeRequest is a struct that represents a request to get a range of members from a sorted set in the cache.
// It contains the key of the sorted set, the start and end of the range, and an optional custom prefix.
type ZRangeRequest struct {
	Key    string
	Start  int64
	End    int64
	Prefix *string
}

// ZRange is a method of Cache that gets a range of members from a sorted set in the cache.
// It takes a context and a pointer to a ZRangeRequest struct,
// and returns the members as a slice of strings and an error.
// The method uses the Redis ZRANGE command to get the members.
func (c *Cache) ZRange(ctx context.Context, request *ZRangeRequest) ([]string, error) {
	return c.client.ZRange(ctx, c.prefixKey(request.Key, request.Prefix), request.Start, request.End).Result()
}

// ZAddRequest is a struct that represents a request to add members to a sorted set in the cache.
// It contains the key of the sorted set, the members to add, and an optional custom prefix.
type ZAddRequest struct {
	Key     string
	Members []*redis.Z
	Prefix  *string
}

// ZAdd is a method of Cache that adds members to a sorted set in the cache.
// It takes a context and a pointer to a ZAddRequest struct,
// and returns the number of members that were added and an error.
// The method uses the Redis ZADD command to add the members.
func (c *Cache) ZAdd(ctx context.Context, request *ZAddRequest) (int64, error) {
	return c.client.ZAdd(ctx, c.prefixKey(request.Key, request.Prefix), request.Members...).Result()
}

// ZScanRequest is a struct that represents a request to incrementally iterate over a sorted set in the cache.
// It contains the key of the sorted set, a cursor to resume the iteration, a match pattern to filter the members,
// a count to limit the number of returned members per call, and an optional custom prefix.
type ZScanRequest struct {
	Key    string
	Cursor uint64
	Match  string
	Count  int64
	Prefix *string
}

// ZScan is a method of Cache that incrementally iterates over a sorted set in the cache.
// It takes a context and a pointer to a ZScanRequest struct,
// and returns the members as a slice of strings, the next cursor, and an error.
// The method uses the Redis ZSCAN command to iterate over the members.
func (c *Cache) ZScan(ctx context.Context, request *ZScanRequest) ([]string, uint64, error) {
	return c.client.
		ZScan(ctx, c.prefixKey(request.Key, request.Prefix), request.Cursor, request.Match, request.Count).
		Result()
}

// ZCardRequest is a struct that represents a request to get the number of members in a sorted set in the cache.
// It contains the key of the sorted set and an optional custom prefix.
type ZCardRequest struct {
	Key    string
	Prefix *string
}

// ZCard is a method of Cache that gets the number of members in a sorted set in the cache.
// It takes a context and a pointer to a ZCardRequest struct,
// and returns the number of members as an int64 and an error.
// The method uses the Redis ZCARD command to get the number of members.
func (c *Cache) ZCard(ctx context.Context, request *ZCardRequest) (int64, error) {
	return c.client.ZCard(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// ZRemRequest is a struct that represents a request to remove members from a sorted set in the cache.
// It contains the key of the sorted set, the members to remove, and an optional custom prefix.
type ZRemRequest struct {
	Key     string
	Members []interface{}
	Prefix  *string
}

// ZRem is a method of Cache that removes members from a sorted set in the cache.
// It takes a context and a pointer to a ZRemRequest struct,
// and returns the number of members that were removed and an error.
// The method uses the Redis ZREM command to remove the members.
func (c *Cache) ZRem(ctx context.Context, request *ZRemRequest) (int64, error) {
	return c.client.ZRem(ctx, c.prefixKey(request.Key, request.Prefix), request.Members).Result()
}

// ZRankRequest is a struct that represents a request to get the rank of a member in a sorted set in the cache.
// It contains the key of the sorted set, the member to find, and an optional custom prefix.
type ZRankRequest struct {
	Key    string
	Member string
	Prefix *string
}

// ZRank is a method of Cache that gets the rank of a member in a sorted set in the cache.
// It takes a context and a pointer to a ZRankRequest struct,
// and returns the rank of the member as an int64 and an error.
// The method uses the Redis ZRANK command to get the rank.
func (c *Cache) ZRank(ctx context.Context, request *ZRankRequest) (int64, error) {
	return c.client.ZRank(ctx, c.prefixKey(request.Key, request.Prefix), request.Member).Result()
}

// SAddRequest is a struct that represents a request to add members to a set in the cache.
// It contains the key of the set, the members to add, and an optional custom prefix.
type SAddRequest struct {
	Key    string
	Value  []interface{}
	Prefix *string
}

// SAdd is a method of Cache that adds members to a set in the cache.
// It takes a context and a pointer to a SAddRequest struct,
// and returns the number of members that were added and an error.
// The method uses the Redis SADD command to add the members.
func (c *Cache) SAdd(ctx context.Context, request *SAddRequest) (int64, error) {
	return c.client.SAdd(ctx, c.prefixKey(request.Key, request.Prefix), request.Value).Result()
}

// SCardRequest is a struct that represents a request to get the number of members in a set in the cache.
// It contains the key of the set and an optional custom prefix.
type SCardRequest struct {
	Key    string
	Prefix *string
}

// SCard is a method of Cache that gets the number of members in a set in the cache.
// It takes a context and a pointer to a SCardRequest struct,
// and returns the number of members as an int64 and an error.
// The method uses the Redis SCARD command to get the number of members.
func (c *Cache) SCard(ctx context.Context, request *SCardRequest) (int64, error) {
	return c.client.SCard(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

// SetNXRequest is a struct that represents a request to set a key in the cache, only if the key does not exist.
// It contains the key, the value to set, the expiration time in seconds, and an optional custom prefix.
type SetNXRequest struct {
	Key     string
	Value   interface{}
	Seconds int64
	Prefix  *string
}

// SetNX is a method of Cache that sets a key in the cache, only if the key does not exist.
// It takes a context and a pointer to a SetNXRequest struct,
// and returns a boolean indicating whether the key was set and an error.
// The method uses the Redis SETNX command to set the key.
func (c *Cache) SetNX(ctx context.Context, request *SetNXRequest) (bool, error) {
	return c.client.SetNX(ctx, c.prefixKey(request.Key, request.Prefix), request.Value, time.Second*time.Duration(request.Seconds)).Result()
}

// LRemRequest is a struct that represents a request to remove occurrences of a value from a list in the cache.
// It contains the key of the list, the count of occurrences to remove, the value to remove, and an optional custom prefix.
type LRemRequest struct {
	Key    string
	Count  int64
	Value  interface{}
	Prefix *string
}

// LRem is a method of Cache that removes occurrences of a value from a list in the cache.
// It takes a context and a pointer to a LRemRequest struct,
// and returns the number of occurrences that were removed and an error.
// The method uses the Redis LREM command to remove the occurrences.
func (c *Cache) LRem(ctx context.Context, request *LRemRequest) (int64, error) {
	return c.client.LRem(ctx, c.prefixKey(request.Key, request.Prefix), request.Count, request.Value).Result()
}
