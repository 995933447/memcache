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

func (s *MemCache) Get(ctx context.Context, req *memcache.GetReq) (*memcache.GetResp, error) {
	var resp memcache.GetResp

	start := time.Now()
	defer func() {
		fastlog.ReportStat("cacheRead", 0, time.Since(start))
	}()

	if req.Key == "" {
		return nil, grpc.NewRPCErrWithMsg(memcache.ErrCode_ErrCodeParamInvalid, "key is required")
	}

	hash, err := util.HashFromCtx(ctx)
	if err != nil {
		fastlog.Warnf("get hash from context err: %v", err)
		return nil, err
	}

	value, ok := cache.Get(req.Key, hash)
	if !ok {
		return &resp, nil
	}

	resp.Ok = true
	resp.Value = value

	return &resp, nil
}
