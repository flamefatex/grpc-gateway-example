package gin

import (
	"context"

	"github.com/flamefatex/grpc-gateway-example/handler"
	"github.com/gin-gonic/gin"
)

// gin
func GetGinRouter(ctx context.Context) *gin.Engine {
	// handler 初始化
	ginRouter := gin.Default()
	handler.ExecRegisterGinHandler(ctx, ginRouter)
	return ginRouter
}
