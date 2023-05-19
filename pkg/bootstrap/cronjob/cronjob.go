package cronjob

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/cronjob"
	"github.com/flamefatex/grpc-gateway-example/definition"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/log"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

// BootstrapCronJob 引导定时任务
func BootstrapCronJob(ctx context.Context) {
	regs := cronjob.ExecRegisterCronJob(ctx)
	for _, reg := range regs {
		// 注入logger
		ctx = ctxzap.ToContext(ctx, definition.OrigZap)
		err := reg.Run(ctx)
		if err != nil {
			log.Errorf("cronjob: %s, run err: %s", reg.Name(), err)
		}
	}
}
