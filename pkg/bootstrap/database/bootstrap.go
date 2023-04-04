package database

import (
	"context"
	"time"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	gormopentracing "gorm.io/plugin/opentracing"
)

func BootstrapDatabase(ctx context.Context) {
	// new db
	dialector := mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        config.Config().GetString("mysql.dsn"),
	})
	logMode := logger.Silent
	if config.Config().GetBool("mysql.enableLog") {
		logMode = logger.Info
	}
	gormConfig := &gorm.Config{
		Logger: logger.Recorder.LogMode(logMode),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		CreateBatchSize: 1000,
	}
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		panic(err)
	}

	// Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// opentracing
	db.Use(gormopentracing.New())
}
