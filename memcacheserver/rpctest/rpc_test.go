package rpctest

import (
	"context"
	"log"
	"testing"

	"github.com/995933447/easymicro/grpc"
	"github.com/995933447/memcache/memcache"
	"github.com/995933447/memcache/memcacheserver/boot"
	"github.com/995933447/memcache/memcacheserver/config"
	"github.com/995933447/memcache/memcacheserver/event"
	"github.com/995933447/runtimeutil"
)

func TestGetValue(t *testing.T) {
	InitEnv()
	
	key := "foofghjktyu34567889"
	err := memcache.SetValue(context.Background(), key, "dddpoiuytyubar2222")
	if err != nil {
		t.Fatal(err)
	}

	value, ok, err := memcache.GetValue(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("not ok")
	}

	t.Log(value)

	err = memcache.DeleteValue(context.Background(), key)
	if err != nil {
		t.Fatal(err)
	}

	value, ok, err = memcache.GetValue(context.Background(), "foo")
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("not ok")
	}

	t.Log(value)
}

func InitEnv() {
	if err := boot.InitNode("memcache"); err != nil {
		log.Fatal(runtimeutil.NewStackErr(err))
	}

	if err := config.LoadConfig(); err != nil {
		log.Fatal(runtimeutil.NewStackErr(err))
	}

	if err := boot.InitMgorm(); err != nil {
		log.Fatal(runtimeutil.NewStackErr(err))
	}

	if err := event.RegisterEventListeners(); err != nil {
		log.Fatal(runtimeutil.NewStackErr(err))
	}

	config.SafeReadServerConfig(func(c *config.ServerConfig) {
		if !c.IsProd() {
			if err := boot.RegisterNatsRPCRoutes(); err != nil {
				log.Fatal(err)
			}
		}
	})

	boot.RegisterGRPCDialOpts()

	config.SafeReadServerConfig(func(c *config.ServerConfig) {
		if err := grpc.PrepareDiscoverGRPC(context.TODO(), memcache.EasymicroGRPCSchema, c.GetDiscoveryName()); err != nil {
			log.Fatal(runtimeutil.NewStackErr(err))
		}

		if err := memcache.PrepareGRPC(c.GetDiscoveryName()); err != nil {
			log.Fatal(runtimeutil.NewStackErr(err))
		}
	})
}
