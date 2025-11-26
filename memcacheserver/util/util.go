package util

import (
	"context"
	"hash/fnv"

	"github.com/995933447/easymicro/grpc"
	"github.com/995933447/memcache/memcache"
	"google.golang.org/grpc/metadata"
)

func Hash(key string) (uint64, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(key))
	if err != nil {
		return 0, err
	}

	hash := uint64(h.Sum32())
	return hash, nil
}

func HashFromCtx(ctx context.Context) (uint64, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, grpc.NewRPCErr(memcache.ErrCode_ErrCodeCtxKeyHashKeyNotFound)
	}

	hashes := meta.Get(grpc.CtxKeyRPCHashKey)
	if len(hashes) == 0 {
		return 0, grpc.NewRPCErr(memcache.ErrCode_ErrCodeCtxKeyHashKeyNotFound)
	}

	hashKey := hashes[0]
	hash, err := Hash(hashKey)
	if err != nil {
		return 0, err
	}
	return hash, nil
}
