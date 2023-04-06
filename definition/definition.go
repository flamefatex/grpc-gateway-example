package definition

import (
	lib_redis "github.com/flamefatex/grpc-gateway-example/pkg/lib/database/redis"
	"go.uber.org/zap"
)

var (
	ServiceName = "grpc-gateway-example"
	GrpcAddr    = ":8081"   // grpc 地址
	HttpAddr    = ":8082"   // http 地址
	OrigZap     *zap.Logger // 原始zap日志对象
	RedisClient lib_redis.RedisClient
)

func SetOrigZap(origZap *zap.Logger) {
	OrigZap = origZap
}

func SetRedisClient(rc lib_redis.RedisClient) {
	RedisClient = rc
}
