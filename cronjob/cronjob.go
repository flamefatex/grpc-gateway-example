package cronjob

import (
	"context"
)

type CronJob interface {
	Name() string
	Run(ctx context.Context) error
}

// ExecRegisterCronJob 注册定时任务
func ExecRegisterCronJob(ctx context.Context) []CronJob {
	regs := []CronJob{}
	return regs
}
