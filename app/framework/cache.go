package framework

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheManager struct {
	Client *redis.Client
}

func (c *CacheManager) Connect(config *Config) {
	c.Client = redis.NewClient(&redis.Options{
		Addr:     config.Cache.Host + ":" + config.Cache.Port,
		DB:       config.Cache.Database,
		Password: config.Cache.Password,
		Username: config.Cache.Username,
	})
}

func (c *CacheManager) Delete(keys ...string) (int64, error) {
	return c.Client.Del(c.Client.Context(), keys...).Result()
}

func (c *CacheManager) Get(key string) (string, error) {
	v, err := c.Client.Get(c.Client.Context(), key).Result()
	if err == redis.Nil {
		err = nil
	}
	return v, err
}

func (c *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	return c.Client.Set(c.Client.Context(), key, value, ttl).Err()
}
