package cache

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Has(key string) (bool, error)
	Get(key string) (interface{}, error)

	// Set takes in a key, value, and expiry in the form of int
	Set(key string, value interface{}, expiry ...int) error
	Forget(key string) error
	// EmptyByMatch forgets everything in the cache based on a pattern
	EmptyByMatch(string) error
	// Empty emptie the entire cache
	Empty() error
}

type RedisCache struct {
	// Connection to Redis
	Conn *redis.Pool
	// Prefix adds unique prefixes to keys to prevent deletion of
	// duplicate keys from multiple applications.
	Prefix string
}

// Entry is the type that is put in the cache to hold values to be serialized
type Entry map[string]interface{}

// Check if a key exists in the redis cache
func (c *RedisCache) Has(str string) (bool, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, err
	}

	return ok, nil
}
