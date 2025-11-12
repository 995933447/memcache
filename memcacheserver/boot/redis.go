package boot

import (
	easymicrorouteredis "github.com/995933447/easymicro/routeredis"
	"github.com/995933447/routeredis"
)

func InitRouteredis() {
	routeredis.OnCmdDone = easymicrorouteredis.FastlogRedisCmd
}
