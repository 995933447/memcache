package confighub

import (
	"github.com/995933447/mconfigcenter/configcenter"
	"github.com/995933447/memcache/memcacheserver/config"
)

type Memcache struct {
	MaxCacheBytes        uint64 `mapstructure:"max_cache_bytes"`
	ForbidInsertAfterOOM bool   `mapstructure:"forbid_insert_after_oom"`
}

type MemcacheConfig struct {
	configcenter.KVSubConfigWrapper[Memcache]
}

func RegisterMemcacheConfig() {
	cfg := &MemcacheConfig{}
	config.SafeReadServerConfig(func(c *config.ServerConfig) {
		if c.MconfigcenterConfigKey != "" {
			cfg.Key = c.MconfigcenterConfigKey
			return
		}
		cfg.Key = "memcache"
	})
	configcenter.RegisterKVSubConfig(cfg.Key, cfg)
}
