package cache

import (
	"sync"
	"time"

	"github.com/995933447/fastlog"
	"github.com/995933447/memcache/memcacheserver/config"
	"github.com/995933447/runtimeutil"
)

var isForbiddenInsert = false

var (
	Hot             = &sync.Map{} // 当前周期内的热key
	LastHot         = &sync.Map{} // 上一轮周期的热key
	GtMyHashHot     = &sync.Map{} // 当前周期内因为key的hash比最后一个节点更靠后而落在第一个节点的热key
	GtMyHashLastHot = &sync.Map{} // 上一轮周期因为key的hash比最后一个节点更靠后而落在第一个节点的热key
)

func init() {
	runtimeutil.Go(func() {
		tk := time.NewTicker(time.Minute * 30)
		defer tk.Stop()

		// 通过30分钟一个周期不断交替把本轮的热key赋值给上轮热key，无锁地把不活跃的key内存释放。
		// 如果key很活跃，30分钟内被连续访问，每次读取key在本轮热key的存储map找不到时候，会试着从上轮热key的存储map中加载同步到本轮到热key存储map中，
		// 不活跃的key会在下轮交替中随着LastHot被赋值替换回收掉。实现1小时内没有被持续访问的key被回收释放内存。
		for {
			<-tk.C
			gcNoHotKey()
			gcGtMyHashNoHotKey()
		}
	})

	runtimeutil.Go(func() {
		tk := time.NewTicker(time.Minute * 3)
		defer tk.Stop()

		for {
			<-tk.C

			var (
				maxCacheBytes        uint64
				forbidInsertAfterOOM bool
			)
			config.SafeReadServerConfig(func(c *config.ServerConfig) {
				maxCacheBytes = c.MaxCacheBytes
				forbidInsertAfterOOM = c.ForbidInsertAfterOOM
			})

			if maxCacheBytes > 0 {
				if cachedBytes := calcCachedBytes(); cachedBytes >= maxCacheBytes {
					fastlog.Warnf("memcache reached max cache size of %d bytes, current size:%d bytes", maxCacheBytes, cachedBytes)
					fastlog.Bill("OOM", "current size:%d bytes over max size %d bytes", cachedBytes, maxCacheBytes)
					gcNoHotKey()
					gcGtMyHashNoHotKey()
					time.Sleep(time.Second * 10)
					if forbidInsertAfterOOM {
						if cachedBytes = calcCachedBytes(); cachedBytes >= maxCacheBytes {
							fastlog.Warnf("memcache cached %d bytes reached max cache size of %d bytes, stop insert new key", cachedBytes, maxCacheBytes)
							fastlog.Bill("stopInsertCazOOM", "current size:%d bytes over max size %d bytes", cachedBytes, maxCacheBytes)
							isForbiddenInsert = true
							continue
						}
					}
				} else {
					fastlog.Debugf("cached bytes:%d", cachedBytes)
				}
				isForbiddenInsert = false
			}
		}
	})
}

func calcCachedBytes() uint64 {
	return runtimeutil.GetCurrMemory()
}

func gcNoHotKey() {
	LastHot = Hot
	Hot = &sync.Map{}
}

func gcGtMyHashNoHotKey() {
	GtMyHashLastHot = GtMyHashHot
	GtMyHashHot = &sync.Map{}
}

func getStore(hash uint64) (*sync.Map, *sync.Map) {
	hot := Hot
	lastHot := LastHot
	if hash > myHash {
		hot = GtMyHashHot
		lastHot = GtMyHashLastHot
	}
	return hot, lastHot
}

func Get(key string, hash uint64) (string, bool) {
	if IsForbiddenKey(key, hash) {
		return "", false
	}

	hot, lastHot := getStore(hash)
	value, ok := hot.Load(key)
	if !ok {
		// 当轮热key中没有，尝试从上一轮的热key中加载
		value, ok = lastHot.Load(key)
		if ok {
			// 缓存到当轮热key
			value, _ = hot.LoadOrStore(key, value)
		}
	}

	if !ok {
		return "", false
	}

	return value.(string), ok
}

func Set(key string, hash uint64, value string) {
	if IsForbiddenKey(key, hash) {
		return
	}

	_, ok := Get(key, hash)
	if !ok && isForbiddenInsert {
		return
	}

	hot, _ := getStore(hash)
	hot.Store(key, value)
}

func Del(key string, hash uint64) {
	if IsForbiddenKey(key, hash) {
		return
	}

	hot, lastHot := getStore(hash)
	lastHot.Delete(key)
	hot.Delete(key)
}
