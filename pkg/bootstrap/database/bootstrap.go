package database

import (
	"context"
)

func BootstrapDatabase(ctx context.Context) {
	bootstrapMysql(ctx)
	bootstrapRedis(ctx)
}
