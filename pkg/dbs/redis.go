package dbs

import (
	"context"
	"e-commerce/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config utils.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     redisAddr(config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})
}

func redisAddr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func Get(c *redis.Client, key string, value interface{}) error {
	strVal, err := c.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(strVal), value)
	if err != nil {
		return err
	}
	return nil
}

func Set(c *redis.Client, key string, value interface{}) error {
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Set(context.Background(), key, str, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetWithExpirationTime(c *redis.Client, key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.Set(context.Background(), key, data, duration).Err()
	if err != nil {
		return err
	}
	return nil
}
