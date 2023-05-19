package cronjob

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/cronjob/cron1"
)

type CronJob interface {
	Name() string
	Run(ctx context.Context) error
}

// ExecRegisterCronJob 注册定时任务
func ExecRegisterCronJob(ctx context.Context) []CronJob {
	regs := []CronJob{
		cron1.GetInstance(),
	}
	return regs
}
