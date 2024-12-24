package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Close() error
	Exist(hash string) (bool, error)
	Set(hash string) error
}

type cache struct {
	client *redis.Client
}

func New(host, port, pass string) (*cache, error) {
	cache := &cache{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: pass,
			DB:       0,
		}),
	}

	err := cache.client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *cache) Close() error {
	return c.client.Close()
}

func (c *cache) Exist(hash string) (bool, error) {
	_, err := c.client.Get(
		context.Background(),
		hash,
	).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *cache) Set(hash string) error {
	return c.client.Set(
		context.Background(),
		hash,
		"",
		0,
	).Err()
}
