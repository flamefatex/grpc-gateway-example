package recovery

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PrintStackHandlerFuncContext(ctx context.Context, p interface{}) (err error) {
	debug.PrintStack()
	return status.Errorf(codes.Internal, "%s", p)
}
