package handler

import (
	"context"
	"time"

	"github.com/995933447/easymicro/grpc"
	"github.com/995933447/fastlog"
	"github.com/995933447/memcache/memcache"
	"github.com/995933447/memcache/memcacheserver/cache"
	"github.com/995933447/memcache/memcacheserver/util"
)

func (s *MemCache) Del(ctx context.Context, req *memcache.DelReq) (*memcache.DelResp, error) {
	var resp memcache.DelResp

	start := time.Now()
	defer func() {
		fastlog.ReportStat("delCache", 0, time.Since(start))
	}()

	if req.Key == "" {
		return nil, grpc.NewRPCErrWithMsg(memcache.ErrCode_ErrCodeParamInvalid, "key is required")
	}

	hash, err := util.HashFromCtx(ctx)
	if err != nil {
		fastlog.Warnf("get hash from context err: %v", err)
		return nil, err
	}

	fastlog.Bill("memcache_del", "del key:%s, hash:%d, my node hash:%d", req.Key, hash, cache.GetMyNodeHash())
	cache.Del(req.Key, hash)

	return &resp, nil
}
