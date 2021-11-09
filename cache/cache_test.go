package cache

import "testing"

func TestRedisCache_Has(t *testing.T) {
	err := testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	// check cache for non-existent key
	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo found in cache, when it shouldn't be")
	}

	// check cache for existent key
	err = testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	inCache, err = testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo not found in cache, when it should be")
	}
}

func TestRedisCache_Get(t *testing.T) {
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	val, err := testRedisCache.Get("foo")
	if err != nil {
		t.Error(err)
	}

	if val != "bar" {
		t.Error("did not get correct value from cache")
	}
}

func TestRedisCache_Forget(t *testing.T) {
	// set it
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	// forget it
	err = testRedisCache.Forget("foo")
	if err != nil {
		t.Error(err)
	}

	// check it
	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo found in cache but it shouldn't be")
	}
}

func TestRedisCache_Empty(t *testing.T) {
	// set it
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	// try to empty
	err = testRedisCache.Empty()
	if err != nil {
		t.Error(err)
	}

	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("foo found in cache but it shouldn't be")
	}
}

func TestRedisCache_EmptyByMatch(t *testing.T) {
	// set values
	err := testRedisCache.Set("foo", "bar")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("bar", "biz")
	if err != nil {
		t.Error(err)
	}

	err = testRedisCache.Set("baz", "bees")
	if err != nil {
		t.Error(err)
	}

	// try to empty by match
	err = testRedisCache.EmptyByMatch("ba")
	if err != nil {
		t.Error(err)
	}

	// check if still in cache
	inCache, err := testRedisCache.Has("foo")
	if err != nil {
		t.Error(err)
	}

	if !inCache {
		t.Error("foo not found in cache but it should be")
	}

	inCache, err = testRedisCache.Has("bar")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("bar found in cache but it shouldn't be")
	}

	inCache, err = testRedisCache.Has("baz")
	if err != nil {
		t.Error(err)
	}

	if inCache {
		t.Error("baz found in cache but it shouldn't be")
	}
}

func TestEncodeDecode(t *testing.T) {
	// make entry
	entry := Entry{"foo": "bar"}

	// encode entry
	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}

	_, err = decode(string(bytes))
	if err != nil {
		t.Error(err)
	}
}
