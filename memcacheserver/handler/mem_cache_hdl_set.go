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

func (s *MemCache) Set(ctx context.Context, req *memcache.SetReq) (*memcache.SetResp, error) {
	var resp memcache.SetResp

	start := time.Now()
	defer func() {
		fastlog.ReportStat("setCache", 0, time.Since(start))
	}()

	if req.Key == "" {
		return nil, grpc.NewRPCErrWithMsg(memcache.ErrCode_ErrCodeParamInvalid, "key is required")
	}

	hash, err := util.HashFromCtx(ctx)
	if err != nil {
		fastlog.Warnf("get hash from context err: %v", err)
		return nil, err
	}

	fastlog.Bill("memcache_set", "set key:%s, hash:%d, my node hash:%d, value:%s", req.Key, hash, cache.GetMyNodeHash(), req.Value)
	cache.Set(req.Key, hash, req.Value)

	return &resp, nil
}
