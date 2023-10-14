package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	prefix string
	client *redis.Client
}

func NewCache(options *redis.Options, prefix string) (*Cache, error) {
	client := redis.NewClient(options)
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := client.Ping(timeout).Err()
	if err != nil {
		return nil, err
	}
	return &Cache{
		prefix: prefix,
		client: client,
	}, nil
}

func (c *Cache) prefixKey(key string, prefix *string) string {
	if prefix != nil {
		return *prefix + key
	}
	return c.prefix + key
}

type SetCacheRequest struct {
	Key     string
	Value   []byte
	Seconds int64
	Prefix  *string
}

// Set Redis `SET key value [expiration]` command.
func (c *Cache) Set(ctx context.Context, request *SetCacheRequest) error {
	_, err := c.client.Set(ctx, c.prefixKey(request.Key, request.Prefix), request.Value, time.Duration(request.Seconds)*time.Second).Result()
	return err
}

type GetCacheRequest struct {
	Key    string
	Prefix *string
}

// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
func (c *Cache) Get(ctx context.Context, request *GetCacheRequest) (string, error) {
	return c.client.Get(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type IncRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) Inc(ctx context.Context, request *IncRequest) (int64, error) {
	return c.client.Incr(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type IncrByRequest struct {
	Key    string
	Value  int64
	Prefix *string
}

func (c *Cache) IncrBy(ctx context.Context, request *IncrByRequest) (int64, error) {
	return c.client.IncrBy(ctx, c.prefixKey(request.Key, request.Prefix), request.Value).Result()
}

type DeleteRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) Delete(ctx context.Context, request *DeleteRequest) (int64, error) {
	return c.client.Del(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type LPushRequest struct {
	Key    string
	Value  []interface{}
	Prefix *string
}

func (c *Cache) LPush(ctx context.Context, request *LPushRequest) (int64, error) {
	return c.client.LPush(ctx, c.prefixKey(request.Key, request.Prefix), request.Value...).Result()
}

type LPopRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) LPop(ctx context.Context, request *LPopRequest) (string, error) {
	return c.client.LPop(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type RPopRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) RPop(ctx context.Context, request *RPopRequest) (string, error) {
	return c.client.RPop(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type LRangeRequest struct {
	Key    string
	Value  string
	Prefix *string
}

func (c *Cache) LRange(ctx context.Context, request *LRangeRequest) ([]string, error) {
	return c.client.LRange(ctx, c.prefixKey(request.Key, request.Prefix), 0, -1).Result()
}

type GetListIndexRequest struct {
	Key    string
	Value  string
	Prefix *string
}

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

type TTLRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) TTL(ctx context.Context, request *TTLRequest) (time.Duration, error) {
	return c.client.TTL(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type ExistsRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) Exists(ctx context.Context, request *ExistsRequest) (bool, error) {
	val, err := c.client.Exists(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

type ExpireRequest struct {
	Key     string
	Seconds int64
	Prefix  *string
}

func (c *Cache) Expire(ctx context.Context, request *ExpireRequest) (bool, error) {
	return c.client.Expire(ctx, c.prefixKey(request.Key, request.Prefix), time.Duration(request.Seconds)*time.Second).Result()
}

type ZRangeRequest struct {
	Key    string
	Start  int64
	End    int64
	Prefix *string
}

func (c *Cache) ZRange(ctx context.Context, request *ZRangeRequest) ([]string, error) {
	return c.client.ZRange(ctx, c.prefixKey(request.Key, request.Prefix), request.Start, request.End).Result()
}

type ZAddRequest struct {
	Key     string
	Members []*redis.Z
	Prefix  *string
}

func (c *Cache) ZAdd(ctx context.Context, request *ZAddRequest) (int64, error) {
	return c.client.ZAdd(ctx, c.prefixKey(request.Key, request.Prefix), request.Members...).Result()
}

type ZScanRequest struct {
	Key    string
	Cursor uint64
	Match  string
	Count  int64
	Prefix *string
}

func (c *Cache) ZScan(ctx context.Context, request *ZScanRequest) ([]string, uint64, error) {
	return c.client.
		ZScan(ctx, c.prefixKey(request.Key, request.Prefix), request.Cursor, request.Match, request.Count).
		Result()
}

type ZCardRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) ZCard(ctx context.Context, request *ZCardRequest) (int64, error) {
	return c.client.ZCard(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type ZRemRequest struct {
	Key     string
	Members []interface{}
	Prefix  *string
}

func (c *Cache) ZRem(ctx context.Context, request *ZRemRequest) (int64, error) {
	return c.client.ZRem(ctx, c.prefixKey(request.Key, request.Prefix), request.Members).Result()
}

type ZRankRequest struct {
	Key    string
	Member string
	Prefix *string
}

func (c *Cache) ZRank(ctx context.Context, request *ZRankRequest) (int64, error) {
	return c.client.ZRank(ctx, c.prefixKey(request.Key, request.Prefix), request.Member).Result()
}

type SAddRequest struct {
	Key    string
	Value  []interface{}
	Prefix *string
}

func (c *Cache) SAdd(ctx context.Context, request *SAddRequest) (int64, error) {
	return c.client.SAdd(ctx, c.prefixKey(request.Key, request.Prefix), request.Value).Result()
}

type SCardRequest struct {
	Key    string
	Prefix *string
}

func (c *Cache) SCard(ctx context.Context, request *SCardRequest) (int64, error) {
	return c.client.SCard(ctx, c.prefixKey(request.Key, request.Prefix)).Result()
}

type SetNXRequest struct {
	Key     string
	Value   interface{}
	Seconds int64
	Prefix  *string
}

func (c *Cache) SetNX(ctx context.Context, request *SetNXRequest) (bool, error) {
	return c.client.SetNX(ctx, c.prefixKey(request.Key, request.Prefix), request.Value, time.Second*time.Duration(request.Seconds)).Result()
}
