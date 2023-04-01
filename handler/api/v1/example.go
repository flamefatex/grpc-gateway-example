package v1

import (
	"context"
	"time"

	"github.com/flamefatex/grpc-gateway-example/pkg/lib/logx"
	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	proto_enum "github.com/flamefatex/grpc-gateway-example/proto/gen/go/enumeration"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type exampleHandler struct {
	proto_v1_example.UnimplementedExampleServiceServer
}

func NewExampleHandler() *exampleHandler {
	return &exampleHandler{}
}

func (h *exampleHandler) List(ctx context.Context, req *proto_v1_example.ExampleListRequest) (resp *proto_v1_example.ExampleListResponse, err error) {
	return nil, nil
}

func (h *exampleHandler) Get(ctx context.Context, req *proto_v1_example.ExampleGetRequest) (resp *proto_v1_example.ExampleGetResponse, err error) {

	now := time.Now()
	resp = &proto_v1_example.ExampleGetResponse{
		Example: &proto_v1_example.Example{
			Id:              1,
			Uuid:            "example-xxx",
			Name:            "示例1",
			Type:            proto_enum.ExampleType_EXAMPLE_TYPE_ONE,
			Description:     "示例1描述",
			CreateTime:      now.Unix(),
			UpdateTime:      now.Unix(),
			CreateTimestamp: timestamppb.New(now),
			UpdateTimestamp: timestamppb.New(now),
		},
	}

	logx.Debug(ctx, "test logx")

	return resp, nil
}

func (h *exampleHandler) Create(ctx context.Context, req *proto_v1_example.ExampleCreateRequest) (resp *proto_v1_example.ExampleCreateResponse, err error) {
	return nil, nil
}

func (h *exampleHandler) Update(ctx context.Context, req *proto_v1_example.ExampleUpdateRequest) (resp *proto_v1_example.ExampleUpdateResponse, err error) {
	return nil, nil
}

func (h *exampleHandler) Delete(ctx context.Context, req *proto_v1_example.ExampleDeleteRequest) (resp *proto_v1_example.ExampleDeleteResponse, err error) {
	return nil, nil
}
