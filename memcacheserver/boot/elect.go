package boot

import (
	"context"

	"github.com/995933447/easymicro/elect"
)

func InitElect(ctx context.Context) error {
	err := elect.InitElect(ctx, &elect.Options{Driver: elect.DriverRedis, RedisConnName: "default"})
	if err != nil {
		return err
	}

	return nil
}
