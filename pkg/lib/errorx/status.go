package errorx

import (
	proto_status "github.com/flamefatex/grpc-gateway-example/proto/gen/go/common/status"
	"google.golang.org/grpc/codes"
)

func OK() *proto_status.Status {
	return &proto_status.Status{
		Code:    int32(FromGRPCCode(codes.OK)),
		Reason:  codes.OK.String(),
		Message: codes.OK.String(),
	}
}
