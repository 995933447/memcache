package memcache

import (
	"github.com/995933447/natsevent"
)

const (
	EventNameMemcacheNodeUp = "memcache.nodeUp" //更新用户信息
)

type MemcacheNodeUpEvent struct {
	Addr string
	Hash uint64
}

func (e *MemcacheNodeUpEvent) Send() error {
	return natsevent.Publish(EventNameMemcacheNodeUp, e)
}
