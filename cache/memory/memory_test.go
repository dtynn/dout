package memory

import (
	"testing"
	"time"

	"github.com/dtynn/dout/cache"
)

func TestSetAndGet(t *testing.T) {
	key := "abcde"
	val := "12345"

	m := New()
	m.Set(key, val)

	v, _ := m.Get(key)
	if v != val {
		t.Errorf("get %s while val is %s", v, val)
	}
}

func TestSetexAndGet(t *testing.T) {
	key := "abcde"
	val := "12345"
	var expire int64 = 2

	m := New()
	m.Setex(key, val, expire)

	v, _ := m.Get(key)
	if v != val {
		t.Errorf("get %s while val is %s", v, val)
	}

	time.Sleep(time.Duration(expire+1) * time.Second)
	v, err := m.Get(key)
	if err != cache.CacheNotFound {
		t.Errorf("get %s which was expected to be expired", v)
	}
}

func TestDelAndGet(t *testing.T) {
	key := "abcde"
	val := "12345"

	m := New()
	m.Set(key, val)

	v, err := m.Get(key)
	if v != val {
		t.Errorf("get %s while val is %s", v, val)
	}
	if err != nil {
		t.Errorf("get err %s which was expected nil", err)
	}

	m.Del(key)
	v, err = m.Get(key)
	if err != cache.CacheNotFound {
		t.Errorf("get %s which was expected to be expired", v)
	}
}
