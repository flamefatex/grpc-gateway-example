package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flamefatex/grpc-gateway-example/model"
	"github.com/flamefatex/grpc-gateway-example/model/query"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/statusx"
	util_paging "github.com/flamefatex/grpc-gateway-example/pkg/util/paging"
	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	proto_enum "github.com/flamefatex/grpc-gateway-example/proto/gen/go/enumeration"
	"github.com/rs/xid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gen"
)

type exampleHandler struct {
	proto_v1_example.UnimplementedExampleServiceServer
}

func NewExampleHandler() *exampleHandler {
	return &exampleHandler{}
}

func (h *exampleHandler) List(ctx context.Context, req *proto_v1_example.ExampleListRequest) (resp *proto_v1_example.ExampleListResponse, err error) {
	resp = &proto_v1_example.ExampleListResponse{}
	req.Paging = util_paging.Normalize(req.Paging)

	q := query.Example
	// 组装条件
	conditions := make([]gen.Condition, 0)
	if req.Uuid != "" {
		conditions = append(conditions, q.Uuid.Eq(req.Uuid))
	}
	if req.Name != "" {
		conditions = append(conditions, q.Name.Like("%"+req.Name+"%"))
	}

	examples, total, err := q.WithContext(ctx).
		Where(conditions...).
		Order(q.Id.Desc()).
		FindByPage(util_paging.OffsetLimit(req.Paging))
	if err != nil {
		err = statusx.Errorf(codes.Internal, "get example list failed, err: %w", err)
		return
	}

	rsExamples := make([]*proto_v1_example.Example, 0)
	for _, example := range examples {
		e := &proto_v1_example.Example{
			Id:              example.Id,
			Uuid:            example.Uuid,
			Name:            example.Name,
			Type:            example.Type,
			Description:     example.Description,
			CreateTime:      example.CreateTime.Unix(),
			UpdateTime:      example.UpdateTime.Unix(),
			CreateTimestamp: timestamppb.New(example.CreateTime),
			UpdateTimestamp: timestamppb.New(example.UpdateTime),
		}
		rsExamples = append(rsExamples, e)
	}

	resp.Examples = rsExamples
	resp.Paging = util_paging.WithTotal(req.Paging, total)
	return
}

func (h *exampleHandler) Get(ctx context.Context, req *proto_v1_example.ExampleGetRequest) (resp *proto_v1_example.ExampleGetResponse, err error) {

	resp = &proto_v1_example.ExampleGetResponse{}

	q := query.Example

	example, err := q.WithContext(ctx).GetByUuid(req.Uuid)
	if err != nil {
		err = statusx.Errorf(codes.Internal, "get example failed, err: %w", err)
		return
	}

	resp.Example = &proto_v1_example.Example{
		Id:              example.Id,
		Uuid:            example.Uuid,
		Name:            example.Name,
		Type:            example.Type,
		Description:     example.Description,
		CreateTime:      example.CreateTime.Unix(),
		UpdateTime:      example.UpdateTime.Unix(),
		CreateTimestamp: timestamppb.New(example.CreateTime),
		UpdateTimestamp: timestamppb.New(example.UpdateTime),
	}

	return
}

func (h *exampleHandler) Create(ctx context.Context, req *proto_v1_example.ExampleCreateRequest) (resp *proto_v1_example.ExampleCreateResponse, err error) {
	resp = &proto_v1_example.ExampleCreateResponse{}

	q := query.Example
	example := &model.Example{
		Uuid:        fmt.Sprintf("example-%s", xid.New().String()),
		Name:        strings.TrimSpace(req.Example.Name),
		Type:        req.Example.Type,
		Description: strings.TrimSpace(req.Example.Description),
	}

	err = q.WithContext(ctx).Create(example)
	if err != nil {
		err = statusx.Errorf(codes.Internal, "create example failed, err: %w", err)
		return
	}

	return
}

func (h *exampleHandler) Update(ctx context.Context, req *proto_v1_example.ExampleUpdateRequest) (resp *proto_v1_example.ExampleUpdateResponse, err error) {
	resp = &proto_v1_example.ExampleUpdateResponse{}

	q := query.Example
	updateParam := map[string]interface{}{
		"name":        strings.TrimSpace(req.Example.Name),
		"description": strings.TrimSpace(req.Example.Description),
	}
	_, err = q.WithContext(ctx).Where(q.Uuid.Eq(req.Example.Uuid)).Updates(updateParam)
	if err != nil {
		err = statusx.Errorf(codes.Internal, "update example failed, err: %w", err)
		return
	}

	return
}

func (h *exampleHandler) Delete(ctx context.Context, req *proto_v1_example.ExampleDeleteRequest) (resp *proto_v1_example.ExampleDeleteResponse, err error) {
	resp = &proto_v1_example.ExampleDeleteResponse{}

	q := query.Example

	_, err = q.WithContext(ctx).DeleteByUuid(req.Uuid)
	if err != nil {
		err = statusx.Errorf(codes.Internal, "delete example failed, err: %w", err)
		return
	}

	return
}

func (h *exampleHandler) Test(ctx context.Context, req *proto_v1_example.ExampleTestRequest) (resp *proto_v1_example.ExampleTestResponse, err error) {
	now := time.Now()
	resp = &proto_v1_example.ExampleTestResponse{
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

	return
}
