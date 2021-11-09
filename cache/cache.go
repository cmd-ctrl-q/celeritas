package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Has(str string) (bool, error)
	Get(str string) (interface{}, error)

	// Set takes in a key, value, and expiry in the form of int
	Set(str string, data interface{}, ttl ...int) error
	Forget(str string) error
	// EmptyByMatch forgets everything in the cache based on a pattern
	EmptyByMatch(str string) error
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

// encode serializes an item before storing it in cache
func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// decode unserializes an item after retrieving it from cache
func decode(str string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))

	d := gob.NewDecoder(&b)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *RedisCache) Get(str string) (interface{}, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	// get item from cache
	cacheEntry, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	// unserialize entry
	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return nil, err
	}

	// get item
	item := decoded[key]

	return item, nil
}

func (c *RedisCache) Set(str string, value interface{}, ttl ...int) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	// store key-value pair in Entry map
	entry := Entry{key: value}

	// serialize entry before storing to redis
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	// add to redis cache
	if len(ttl) > 0 {
		// expiry
		// SETEX is set with expiry
		_, err := conn.Do("SETEX", key, ttl[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		// no expiry
		_, err := conn.Do("SET", key, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Forget(str string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

// EmptyByMatch removes entrie from cache based on a pattern
func (c *RedisCache) EmptyByMatch(str string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Conn.Get()
	defer conn.Close()

	keys, err := c.getkeys(key)
	if err != nil {
		return err
	}

	// forget keys
	for _, x := range keys {
		_, err = conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Empty() error {
	key := fmt.Sprintf("%s:", c.Prefix)
	conn := c.Conn.Get()
	defer conn.Close()

	keys, err := c.getkeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err = conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) getkeys(pattern string) ([]string, error) {
	conn := c.Conn.Get()
	defer conn.Close()

	iter := 0
	var keys []string

	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}

		// iterate to next entry
		iter, _ = redis.Int(arr[0], nil)   // arr[0] holds key
		k, _ := redis.Strings(arr[1], nil) // arr[1] holds value
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
