package cache

import (
	"fmt"
	"time"

	"github.com/995933447/easymicro/node"
	"github.com/995933447/fastlog"
	"github.com/995933447/memcache/memcache"
	"github.com/995933447/memcache/memcacheserver/util"
	"github.com/995933447/natsevent"
	"github.com/995933447/runtimeutil"
)

var myHash uint64

var (
	forbidAnyKeyTimeoutAt     time.Time
	forbidGtMyHashTimeoutAt   time.Time
	forbidLteKeyHashTimeoutAt time.Time
)

var (
	forbidGtMyHashAndLteKeyHash uint64
	forbidLteKeyHash            uint64
)

func OnNodeUp() error {
	forbidAnyKeyTimeoutAt = time.Now().Add(time.Minute)

	myHost, err := node.GetOrAutoSetHost()
	if err != nil {
		return err
	}

	myPort, err := node.GetOrAutoSetPort()
	if err != nil {
		return err
	}

	myAddr := fmt.Sprintf("%s:%d", myHost, myPort)

	myHash, err = util.Hash(myAddr)
	if err != nil {
		return err
	}

	fastlog.EnableStdoutPrinter()
	fastlog.Infof("my node hash: %d", myHash)
	fastlog.DisableStdoutPrinter()

	err = (&memcache.MemcacheNodeUpEvent{Addr: myAddr, Hash: myHash}).Send()
	if err != nil {
		return err
	}

	err = natsevent.Subscribe(memcache.EventNameMemcacheNodeUp, memcache.EasymicroGRPCPbServiceNameMemCache, func(evt *memcache.MemcacheNodeUpEvent) error {
		// 不处理集群中自己节点启动的扩容事件
		if evt.Addr == myAddr {
			return nil
		}

		OnNodeScaled(evt.Addr, evt.Hash)
		return nil
	})
	if err != nil {
		return err
	}

	fastlog.Bill("nodeUp", "im up,my addr:%s, my hash:%d", myAddr, myHash)
	return nil
}

func OnNodeScaled(addr string, hash uint64) {
	fastlog.Bill("scaledNode", "Listen scaled node, addr:%s, hash:%d", addr, hash)

	if hash <= myHash {
		isCleaningKeyCazScaled := forbidLteKeyHashTimeoutAt.After(time.Now())
		if !isCleaningKeyCazScaled || hash >= forbidLteKeyHash {
			forbidLteKeyHash = hash
			forbidLteKeyHashTimeoutAt = time.Now().Add(time.Minute)
			runtimeutil.Go(func() {
				time.Sleep(time.Second)
				gcNoHotKey()
				gcGtMyHashNoHotKey()
				time.Sleep(time.Second * 50)
				gcNoHotKey()
				gcGtMyHashNoHotKey()
			})
		}
		return
	}

	isCleaningKeyCazScaled := forbidLteKeyHashTimeoutAt.After(time.Now())
	if !isCleaningKeyCazScaled || hash >= forbidGtMyHashAndLteKeyHash {
		forbidGtMyHashAndLteKeyHash = hash
		forbidGtMyHashTimeoutAt = time.Now().Add(time.Minute)
		runtimeutil.Go(func() {
			time.Sleep(time.Second)
			gcGtMyHashNoHotKey()
			time.Sleep(time.Second * 50)
			gcGtMyHashNoHotKey()
		})
	}
}

// IsForbiddenKey 根据key hash判断是否禁止进行操作，新节点加入集群时候，新起节点在一定时间内不做读写,和相邻节点要在一定时间内对key进行过滤，
// 避免部分rpc客户端异步发现服务新节点延迟，部分请求仍然写入相邻节点，导致缓存不一致问题。
func IsForbiddenKey(key string, hash uint64) bool {
	if forbidAnyKeyTimeoutAt.After(time.Now()) {
		fastlog.Bill("hitForbidCacheAnyKey", "forbidAnyKey key:%s, hash:%d, my node hash:%d, timeout at:%v", key, hash, myHash, forbidAnyKeyTimeoutAt.String())
		return true
	}

	if hash > myHash {
		if (hash <= forbidGtMyHashAndLteKeyHash && forbidGtMyHashTimeoutAt.After(time.Now())) ||
			forbidLteKeyHashTimeoutAt.After(time.Now()) {
			fastlog.Bill("hitForbidCacheGteMyHash", "forbidGteMyHash key:%s, hash:%d, my node hash:%d, timeout at:%v or %v", key, hash, myHash, forbidGtMyHashTimeoutAt.String(), forbidLteKeyHashTimeoutAt.String())
			return true
		}
	} else if hash <= forbidLteKeyHash && forbidLteKeyHashTimeoutAt.After(time.Now()) {
		fastlog.Bill("hitForbidCacheLteKeyHash", "forbidLteKeyHash key:%s, hash:%d, my node hash:%d, timeout at:%v", key, hash, myHash, forbidLteKeyHashTimeoutAt.String())
		return true
	}

	return false
}

func GetMyNodeHash() uint64 {
	return myHash
}
