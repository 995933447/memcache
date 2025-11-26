package handler

import (
	"github.com/995933447/memcache/memcache"
)

type MemCache struct {
	memcache.UnimplementedMemCacheServer
	ServiceName string
}

var MemCacheHandler = &MemCache{
	ServiceName: memcache.EasymicroGRPCPbServiceNameMemCache,
}
