package boot

import "github.com/995933447/memcache/memcacheserver/cache"

func InitApp() {
	if err := InitConfigHub(); err != nil {
		panic(err)
	}

	if err := cache.OnNodeUp(); err != nil {
		panic(err)
	}
}
