package v1

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/flamefatex/grpc-gateway-example/model"
	"github.com/flamefatex/grpc-gateway-example/model/query"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/errorx"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/logx"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/pagingx"
	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	proto_enum "github.com/flamefatex/grpc-gateway-example/proto/gen/go/enumeration"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gen"
)

type exampleHandler struct {
	proto_v1_example.UnimplementedExampleServiceServer
}

func NewExampleHandler() *exampleHandler {
	return &exampleHandler{}
}

func (h *exampleHandler) List(ctx context.Context, req *proto_v1_example.ListRequest) (resp *proto_v1_example.ListResponse, err error) {
	resp = &proto_v1_example.ListResponse{}
	req.Paging = pagingx.Normalize(req.Paging)

	q := query.Example
	// 组装条件
	conditions := make([]gen.Condition, 0)
	if req.Id != "" {
		conditions = append(conditions, q.Id.Eq(req.Id))
	}
	if req.Name != "" {
		conditions = append(conditions, q.Name.Like("%"+req.Name+"%"))
	}

	examples, total, err := q.WithContext(ctx).
		Where(conditions...).
		Order(q.Id.Desc()).
		FindByPage(pagingx.OffsetLimit(req.Paging))
	if err != nil {
		err = errorx.ErrorfInternalServer("EXAMPLE_LIST_ERROR", "get example list failed, err: %s", err)
		return
	}

	rsExamples := make([]*proto_v1_example.Example, 0)
	for _, example := range examples {
		e := &proto_v1_example.Example{
			Id:          example.Id,
			Name:        example.Name,
			Type:        example.Type,
			Description: example.Description,
			CreatedAt:   timestamppb.New(example.CreatedAt),
			UpdatedAt:   timestamppb.New(example.UpdatedAt),
		}
		rsExamples = append(rsExamples, e)
	}

	resp.Examples = rsExamples
	resp.Paging = pagingx.WithTotal(req.Paging, total)
	return
}

func (h *exampleHandler) Get(ctx context.Context, req *proto_v1_example.GetRequest) (resp *proto_v1_example.GetResponse, err error) {

	resp = &proto_v1_example.GetResponse{}

	q := query.Example

	example, err := q.WithContext(ctx).GetById(req.Id)
	if err != nil {
		err = errorx.ErrorfInternalServer("EXAMPLE_GET_ERROR", "get example failed, err: %s", err)
		return
	}

	resp.Example = &proto_v1_example.Example{
		Id:          example.Id,
		Name:        example.Name,
		Type:        example.Type,
		Description: example.Description,
		CreatedAt:   timestamppb.New(example.CreatedAt),
		UpdatedAt:   timestamppb.New(example.UpdatedAt),
	}

	return
}

func (h *exampleHandler) Create(ctx context.Context, req *proto_v1_example.CreateRequest) (resp *proto_v1_example.CreateResponse, err error) {
	resp = &proto_v1_example.CreateResponse{}

	q := query.Example
	example := &model.Example{
		Id:          fmt.Sprintf("example-%s", xid.New().String()),
		Name:        strings.TrimSpace(req.Example.Name),
		Type:        req.Example.Type,
		Description: strings.TrimSpace(req.Example.Description),
	}

	err = q.WithContext(ctx).Create(example)
	if err != nil {
		err = errorx.ErrorfInternalServer("EXAMPLE_CREATE_ERROR", "create example failed, err: %s", err)
		return
	}

	return
}

func (h *exampleHandler) Update(ctx context.Context, req *proto_v1_example.UpdateRequest) (resp *proto_v1_example.UpdateResponse, err error) {
	resp = &proto_v1_example.UpdateResponse{}

	q := query.Example
	updateParam := map[string]interface{}{
		"name":        strings.TrimSpace(req.Example.Name),
		"description": strings.TrimSpace(req.Example.Description),
	}
	_, err = q.WithContext(ctx).Where(q.Id.Eq(req.Example.Id)).Updates(updateParam)
	if err != nil {
		err = errorx.ErrorfInternalServer("EXAMPLE_UPDATE_ERROR", "update example failed, err: %s", err)
		return
	}

	return
}

func (h *exampleHandler) Delete(ctx context.Context, req *proto_v1_example.DeleteRequest) (resp *proto_v1_example.DeleteResponse, err error) {
	resp = &proto_v1_example.DeleteResponse{}

	q := query.Example

	_, err = q.WithContext(ctx).DeleteById(req.Id)
	if err != nil {
		err = errorx.ErrorfInternalServer("EXAMPLE_DELETE_ERROR", "delete example failed, err: %s", err)
		return
	}

	return
}

func (h *exampleHandler) TestCustomHttp(ctx context.Context, req *proto_v1_example.TestCustomHttpRequest) (resp *empty.Empty, err error) {

	code := "401"
	if req.Code != "" {
		code = req.Code
	}

	appId := ""
	appIds := metadata.ValueFromIncomingContext(ctx, "x-app-id")
	if len(appIds) > 0 {
		appId = appIds[0]
		logx.Infof(ctx, "app-id:%s", appId)
	}

	// http code
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", code))

	// http header
	// Set-Cookie: <cookie-name>=<cookie-value>; Max-Age=<non-zero-digit>
	_ = grpc.SetHeader(ctx, metadata.Pairs(runtime.MetadataPrefix+"Set-Cookie", "mySessionId:xxxx;Max-Age=180"))

	//
	_ = grpc.SetHeader(ctx, metadata.Pairs("x-ffx-token", "xxx.yyy.zzz"))

	return nil, nil
}

func (h *exampleHandler) TestError(ctx context.Context, req *proto_v1_example.TestErrorRequest) (resp *proto_v1_example.TestErrorResponse, err error) {
	now := time.Now()
	resp = &proto_v1_example.TestErrorResponse{
		Example: &proto_v1_example.Example{
			Id:          "example-xxx",
			Name:        "示例1",
			Type:        proto_enum.ExampleType_EXAMPLE_TYPE_ONE,
			Description: "示例1描述",
			CreatedAt:   timestamppb.New(now),
			UpdatedAt:   timestamppb.New(now),
		},
	}

	err = errors.New("test error")
	err = errorx.ErrorfInternalServer("EXAMPLE_DELETE_ERROR", "delete example failed, err: %s", err).
		WithCause(err).
		WithMetadata(map[string]string{
			"aa": "aad",
		})
	return
}
