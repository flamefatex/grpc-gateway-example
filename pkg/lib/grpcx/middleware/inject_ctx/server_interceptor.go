package inject_ctx

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx := ctx

		if len(o.injectFromOrigCtxFuncs) != 0 {
			for _, fn := range o.injectFromOrigCtxFuncs {
				newCtx = fn(newCtx)
			}
		}

		return handler(newCtx, req)
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOptions(opts)

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		newCtx := stream.Context()

		if len(o.injectFromOrigCtxFuncs) != 0 {
			for _, fn := range o.injectFromOrigCtxFuncs {
				newCtx = fn(newCtx)
			}
		}

		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = newCtx
		err = handler(srv, wrapped)

		return err
	}
}
