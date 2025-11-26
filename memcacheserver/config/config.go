package config

import (
	"sync"

	"github.com/995933447/easymicro/loader"
	"github.com/995933447/mconfigcenter/configcenter"
	"github.com/995933447/memcache/memcache"
)

const ServerConfigFileName = "memcacheserver"

type ServerConfig struct {
	SamplePProfTimeLongSec       int    `mapstructure:"sample_pprof_time_long_sec"`
	Env                          string `mapstructure:"env"`
	EnabledMconfigcenter         bool   `mapstructure:"enabled_mconfigcenter"`
	MconfigcenterKVConfigDataSrc int    `mapstructure:"mconfigcenter_kv_config_data_src"`
	MconfigcenterListenerGroup   string `mapstructure:"mconfigcenter_listener_group"`
	MconfigcenterConfigKey       string `mapstructure:"mconfigcenter_config_key"`
	MconfigcenterLocalImgMgoConn string `mapstructure:"mconfigcenter_local_img_mgo_conn"`
	MconfigcenterLocalImgMgoDb   string `mapstructure:"mconfigcenter_local_img_mgo_db"`
	MaxCacheBytes                uint64 `mapstructure:"max_cache_bytes"`
	ForbidInsertAfterOOM         bool   `mapstructure:"forbid_insert_after_oom"`
	DiscoveryName                string `mapstructure:"discovery_name"`
}

func (c *ServerConfig) GetDiscoveryName() string {
	if c.DiscoveryName == "" {
		return memcache.EasymicroDiscoveryName
	}
	return c.DiscoveryName
}

func (c *ServerConfig) IsDev() bool {
	return c.Env == "dev"
}

func (c *ServerConfig) IsTest() bool {
	return c.Env == "test"
}

func (c *ServerConfig) IsProd() bool {
	return c.Env == "prod"
}

var (
	serverConfig   ServerConfig
	serverConfigMu sync.RWMutex
)

func getServerConfig() *ServerConfig {
	return &serverConfig
}

func SafeReadServerConfig(fn func(c *ServerConfig)) {
	serverConfigMu.RLock()
	defer serverConfigMu.RUnlock()
	fn(getServerConfig())
}

func SafeWriteServerConfig(fn func(c *ServerConfig)) {
	serverConfigMu.Lock()
	defer serverConfigMu.Unlock()
	fn(getServerConfig())
}

func LoadConfig() error {
	var err error
	err = loader.LoadFastlogFromLocal(nil)
	if err != nil {
		return err
	}

	err = loader.LoadAndWatchConfig(ServerConfigFileName, &serverConfig, &serverConfigMu, nil)
	if err != nil {
		return err
	}

	if err = loader.LoadEtcdFromLocal(); err != nil {
		return err
	}

	if err = loader.LoadDiscoveryFromLocal(); err != nil {
		return err
	}

	if err = loader.LoadNatsFromLocal(); err != nil {
		return err
	}

	var enabledMongo bool
	SafeReadServerConfig(func(c *ServerConfig) {
		enabledMongo = c.EnabledMconfigcenter && c.MconfigcenterKVConfigDataSrc == int(configcenter.KVConfigDataSrcLocalImage)
	})

	if enabledMongo {
		if err = loader.LoadAndWatchMongoFromLocal(); err != nil {
			return err
		}
	}

	return nil
}
