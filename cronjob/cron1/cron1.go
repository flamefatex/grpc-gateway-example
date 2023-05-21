package cron1

import (
	"context"
	"sync"
	"time"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/config"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/logx"
)

var (
	name     = "Cron1"
	instance *Cron1 // 单例实体
	once     sync.Once
)

type Cron1 struct{}

func GetInstance() *Cron1 {
	once.Do(func() {
		instance = &Cron1{}
	})
	return instance
}

func (cj *Cron1) Name() string {
	return name
}

func (cj *Cron1) Run(ctx context.Context) (err error) {

	// 判断是否开启定时任务
	if !config.Config().GetBool("app.cronjob.cron1.enabled") {
		return nil
	}

	// 开始定时任务
	t := time.NewTicker(config.Config().GetDuration("app.cronjob.cron1.interval"))
	go func() {
		for {
			startTime := time.Now()
			log.Infof("cronjob: %s start", cj.Name())

			// do your job

			duration := time.Now().Sub(startTime)
			logx.Infof(ctx, "cronjob: %s end, duration: %s", cj.Name(), duration.String())
			// wait
			<-t.C
		}
	}()

	return nil
}