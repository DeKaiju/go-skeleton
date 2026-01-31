package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var c = cache.New(cache.NoExpiration, cache.NoExpiration)

func Set(k string, v interface{}, d time.Duration) {
	c.Set(k, v, d)
}

func Get(k string) (interface{}, bool) {
	return c.Get(k)
}
