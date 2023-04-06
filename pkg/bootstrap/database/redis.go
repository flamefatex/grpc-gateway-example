package database

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	lib_redis "github.com/flamefatex/grpc-gateway-example/pkg/lib/database/redis"
)

func bootstrapRedis(ctx context.Context) {
	conf := &lib_redis.Config{}
	switch config.Config().GetString("redis.mode") {
	case lib_redis.ModeSentinel:
		conf = &lib_redis.Config{
			Mode:           lib_redis.ModeSentinel,
			Password:       config.Config().GetString("redis.password"),
			Db:             config.Config().GetInt("redis.db"),
			ConnectTimeout: config.Config().GetDuration("redis.connectTimeout"),
			ReadTimeout:    config.Config().GetDuration("redis.readTimeout"),
			WriteTimeout:   config.Config().GetDuration("redis.writeTimeout"),
			Sentinel: &lib_redis.Sentinel{
				Addrs:          config.Config().GetStringSlice("redis.sentinel.addrs"),
				MasterName:     config.Config().GetString("redis.sentinel.masterName"),
				Password:       config.Config().GetString("redis.sentinel.password"),
				ConnectTimeout: config.Config().GetDuration("redis.sentinel.connectTimeout"),
				ReadTimeout:    config.Config().GetDuration("redis.sentinel.readTimeout"),
				WriteTimeout:   config.Config().GetDuration("redis.sentinel.writeTimeout"),
			},
		}
	default:
		conf = &lib_redis.Config{
			Mode:           lib_redis.ModeNormalDsn,
			Dsn:            config.Config().GetString("redis.dsn"),
			ConnectTimeout: config.Config().GetDuration("redis.connectTimeout"),
			ReadTimeout:    config.Config().GetDuration("redis.readTimeout"),
			WriteTimeout:   config.Config().GetDuration("redis.writeTimeout"),
		}
	}
	rc, _ := lib_redis.NewRedisClient(conf)

	definition.SetRedisClient(rc)
}
