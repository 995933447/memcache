package memcache

import (
	"context"
	"fmt"

	easymicrogrpc "github.com/995933447/easymicro/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// PrepareGRPC 完成调用memcache服务grpc准备工作
func PrepareGRPC(discoveryName string, dialGRPCOpts ...grpc.DialOption) error {
	if err := easymicrogrpc.PrepareDiscoverGRPC(context.TODO(), EasymicroGRPCSchema, discoveryName); err != nil {
		return err
	}

	dialGRPCOpts = append(dialGRPCOpts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy": "%s"}`, easymicrogrpc.BalancerNameFnvConsistentHash1aSum32)))
	easymicrogrpc.RegisterServiceDialOpts(EasymicroGRPCPbServiceNameMemCache, true, dialGRPCOpts...)

	return nil
}

func GetValue(ctx context.Context, key string, opts ...grpc.CallOption) (string, bool, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, easymicrogrpc.CtxKeyRPCHashKey, key)
	getResp, err := MemCacheGRPC().Get(ctx, &GetReq{
		Key: key,
	}, opts...)
	if err != nil {
		return "", false, err
	}
	return getResp.Value, getResp.Ok, nil
}

func SetValue(ctx context.Context, key string, value string, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, easymicrogrpc.CtxKeyRPCHashKey, key)
	_, err := MemCacheGRPC().Set(ctx, &SetReq{
		Key:   key,
		Value: value,
	}, opts...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteValue(ctx context.Context, key string, opts ...grpc.CallOption) error {
	ctx = metadata.AppendToOutgoingContext(ctx, easymicrogrpc.CtxKeyRPCHashKey, key)
	_, err := MemCacheGRPC().Del(ctx, &DelReq{
		Key: key,
	}, opts...)
	if err != nil {
		return err
	}
	return nil
}
