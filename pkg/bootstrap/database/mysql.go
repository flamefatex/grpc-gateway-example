package database

import (
	"context"
	"time"

	"github.com/flamefatex/grpc-gateway-example/model/query"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	gormopentracing "gorm.io/plugin/opentracing"
)

func bootstrapMysql(ctx context.Context) {
	if !config.Config().GetBool("mysql.enabled") {
		return
	}

	// new db
	dialector := mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        config.Config().GetString("mysql.dsn"),
	})
	logMode := logger.Silent
	if config.Config().GetBool("mysql.logEnabled") {
		logMode = logger.Info
	}
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		CreateBatchSize: 1000,
	}
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {

		panic("sssssss")
		// log.Fatalf("open mysql database failed, err:%s", err)
	}

	// Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("open mysql database connection pool failed, err:%s", err)

		panic(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// opentracing
	if config.Config().GetBool("mysql.opentracingEnabled") {
		db.Use(gormopentracing.New())
	}

	// gorm prometheus
	// db.Use(prometheus.New(prometheus.Config{
	//	DBName:          "example", // `DBName` as metrics label
	//	RefreshInterval: 15,        // refresh metrics interval (default 15 seconds)
	//	MetricsCollector: []prometheus.MetricsCollector{
	//		&prometheus.MySQL{VariableNames: []string{"Threads_running"}},
	//	},
	// }))

	// bind default db instance
	query.SetDefault(db)
}
