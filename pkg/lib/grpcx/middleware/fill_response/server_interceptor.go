package fill_response

import (
	"context"

	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		// 处理后发送前，进行填充
		if o.fillResponseFunc != nil {
			resp = o.fillResponseFunc(ctx, resp, err)
		}
		return resp, err
	}
}
