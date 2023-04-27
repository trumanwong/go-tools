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

func NewCache(addr, password, prefix string, db int) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
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

func (c *Cache) prefixKey(key string) string {
	return c.prefix + key
}

// Set Redis `SET key value [expiration]` command.
func (c *Cache) Set(ctx context.Context, key string, value []byte, seconds int64) error {
	_, err := c.client.Set(ctx, c.prefixKey(key), value, time.Duration(seconds)*time.Second).Result()
	return err
}

// Get Redis `GET key` command. It returns redis.Nil error when key does not exist.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) Inc(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, c.prefixKey(key), value).Result()
}

func (c *Cache) Delete(ctx context.Context, key string) (int64, error) {
	return c.client.Del(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) LPush(ctx context.Context, key string, value ...interface{}) (int64, error) {
	return c.client.LPush(ctx, c.prefixKey(key), value).Result()
}

func (c *Cache) LPop(ctx context.Context, key string) (string, error) {
	return c.client.LPop(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) RPop(ctx context.Context, key string) (string, error) {
	return c.client.RPop(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) GetListIndex(ctx context.Context, key string, value string) (int, error) {
	list, err := c.client.LRange(ctx, c.prefixKey(key), 0, -1).Result()
	if err != nil {
		return 0, err
	}
	for i, v := range list {
		if v == value {
			return len(list) - i, nil
		}
	}
	return 0, errors.New("no found")
}

func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	val, err := c.client.Exists(ctx, c.prefixKey(key)).Result()
	if err != nil {
		return false, err
	}
	return val == 1, nil
}

func (c *Cache) Expire(ctx context.Context, key string, seconds int64) (bool, error) {
	return c.client.Expire(ctx, c.prefixKey(key), time.Duration(seconds)*time.Second).Result()
}

func (c *Cache) ZRange(ctx context.Context, key string, start, end int64) ([]string, error) {
	return c.client.ZRange(ctx, c.prefixKey(key), start, end).Result()
}

type Z struct {
	Score  float64
	Member interface{}
}

func (c *Cache) ZAdd(ctx context.Context, key string, members ...*Z) (int64, error) {
	arr := make([]*redis.Z, len(members))
	for i, v := range members {
		arr[i] = &redis.Z{
			Score:  v.Score,
			Member: v.Member,
		}
	}
	return c.client.ZAdd(ctx, c.prefixKey(key), arr...).Result()
}

func (c *Cache) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.
		ZScan(ctx, c.prefixKey(key), cursor, match, count).
		Result()
}

func (c *Cache) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, c.prefixKey(key)).Result()
}

func (c *Cache) ZRem(ctx context.Context, key string, members []interface{}) (int64, error) {
	return c.client.ZRem(ctx, c.prefixKey(key), members).Result()
}

func (c *Cache) ZRank(ctx context.Context, key, member string) (int64, error) {
	return c.client.ZRank(ctx, c.prefixKey(key), member).Result()
}
