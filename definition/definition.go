package definition

import "go.uber.org/zap"

var (
	ServiceName = "grpc-gateway-example"
	GrpcAddr    = ":8081"   // grpc 地址
	HttpAddr    = ":8082"   // http 地址
	OrigZap     *zap.Logger // 原始zap日志对象
)

func SetOrigZap(origZap *zap.Logger) {
	OrigZap = origZap
}
