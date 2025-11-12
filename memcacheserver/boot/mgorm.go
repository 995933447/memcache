package boot

import (
	easymciromgorm "github.com/995933447/easymicro/mgorm"

	"github.com/995933447/mgorm"
)

func InitMgorm() error {

	mgorm.OnQueryDone = easymciromgorm.FastlogMgormQuery

	return nil
}
