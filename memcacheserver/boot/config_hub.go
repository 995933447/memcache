package boot

import (
	"github.com/995933447/mconfigcenter/configcenter"
	"github.com/995933447/memcache/memcacheserver/config"
	"github.com/995933447/memcache/memcacheserver/confighub"
)

func InitConfigHub() error {
	var listenerGroup string
	var enabledMconfigcenter bool
	config.SafeReadServerConfig(func(c *config.ServerConfig) {
		listenerGroup = c.MconfigcenterListenerGroup
		enabledMconfigcenter = c.EnabledMconfigcenter
	})
	if !enabledMconfigcenter {
		return nil
	}

	err := configcenter.InitReconfmgrReloader(listenerGroup)
	if err != nil {
		return err
	}

	var mgoConn, mgoDb string
	var kvConfigDataSrc int
	config.SafeReadServerConfig(func(c *config.ServerConfig) {
		mgoConn = c.MconfigcenterLocalImgMgoConn
		mgoDb = c.MconfigcenterLocalImgMgoDb
		kvConfigDataSrc = c.MconfigcenterKVConfigDataSrc
	})

	err = configcenter.RegisterKVConfig(1, nil, mgoConn, mgoDb, configcenter.KVConfigDataSrc(kvConfigDataSrc))
	if err != nil {
		return err
	}

	confighub.RegisterMemcacheConfig()

	return nil
}
