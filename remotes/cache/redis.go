package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joaopandolfi/blackwhale/configurations"
	"github.com/joaopandolfi/blackwhale/utils"
)

var rcache *redisCache

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func GetRedis() Cache {
	cfg := configurations.Configuration.Redis
	return initializeRedis(cfg.Server, cfg.Password, cfg.DB)
}

func initializeRedis(server, password string, db int) Cache {
	if rcache == nil {
		utils.Info("[CACHE] using redis cache", "Redis")

		rcache = &redisCache{
			ctx: context.Background(),
			client: redis.NewClient(&redis.Options{
				Addr:     server,
				Password: password,
				DB:       db,
			}),
		}
	}

	return rcache
}

func (c *redisCache) Put(key string, data interface{}, duration time.Duration) error {
	err := c.client.Set(c.ctx, key, data, duration).Err()
	if err != nil {
		return fmt.Errorf("putting data on key (%s): %w", key, err)
	}
	return nil
}

func (c *redisCache) Get(key string) (interface{}, error) {
	val, err := c.client.Get(c.ctx, "key").Result()
	if err == redis.Nil {
		return nil, nil // Key does not exists
	} else if err != nil {
		return nil, fmt.Errorf("getting data from %s : %w", key, err)
	}
	return val, nil
}

func (c *redisCache) Delete(key string) error {
	c.client.Del(c.ctx, key)
	return nil
}
func (c *redisCache) Size() int {
	return 0
}
func (c *redisCache) Flush() error {
	return nil
}
func (c *redisCache) GracefullShutdown() {
	c.client.Close()
}
