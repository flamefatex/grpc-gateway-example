# grpc-gateway-example

grpc-gateway的使用示例，其中集成了各类库。

# 集成

* grpc-gateway
* viper
* zap
* opentracing
* gorm
* redigo

# 目录结构
```
.
├── cmd                                 // 命令
├── comsumer                            // 消费者
├── cronjob                             // 定时任务
├── definition                          // 全局变量定义
├── hack                                // 自动化相关
├── handler                             // 处理器
├── logic                               // 业务逻辑层
├── model                               // 数据模型、DAO层
├── pkg
│     ├── bootstrap                     // 引导
│     ├── lib                           // 核心库
│     └── util                          // 公共工具库
└── proto
    ├── gen
    └── src                             // proto源码
```

# 使用

## 1.用protocol buffer 定义grpc服务并生成go文件

参考`proto/src/api/v1/example/example.proto`
```protobuf
syntax = "proto3";

package flamefatex.grpc_gateway_example.api.v1.example;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/empty.proto";
import "validate/validate.proto";

import  "common/paging/paging.proto";
import  "common/status/status.proto";
import  "enumeration/example.proto";

// ExampleService 示例服务
service ExampleService {
  // List 查询示例列表
  rpc List(ExampleListRequest) returns (ExampleListResponse) {
    option (google.api.http) = {
      get: "/api/v1/example/list"
    };
  }
}

// Example 示例
message Example {
  // ID
  int64 id = 1;
  // 唯一标识
  string uuid = 2;
  // 名称
  string name = 3;
  // 类型
  flamefatex.grpc_gateway_example.enumeration.ExampleType type = 4;
  // 描述
  string description = 5;
  // 创建时间
  int64 create_time = 6;
  // 更新时间
  int64 update_time = 7;
  // 创建时间
  google.protobuf.Timestamp create_timestamp = 8;
  // 更新时间
  google.protobuf.Timestamp update_timestamp = 9;
}

// ExampleListRequest 查询示例列表请求
message ExampleListRequest {
  // 分页
  flamefatex.grpc_gateway_example.common.paging.Paging paging = 1;
  // uuid
  string uuid = 2;
  // 名称
  string name = 3;
}

// ExampleListResponse 查询示例列表响应
message ExampleListResponse {
  // 请求ID
  string request_id = 1;
  // 请求状态
  flamefatex.grpc_gateway_example.common.status.Status status = 2;
  // 分页
  flamefatex.grpc_gateway_example.common.paging.Paging paging = 3;
  // 示例列表
  repeated Example examples = 4;
}

```

执行makefile指令
```shell
make buf-gen
```

## 2.定义模型与生成query DAO
参考`model/example.go`
```go
package model

import (
	"time"

	proto_enum "github.com/flamefatex/grpc-gateway-example/proto/gen/go/enumeration"
	"gorm.io/gen"
)

type Example struct {
	Id          int64
	Uuid        string
	Name        string
	Type        proto_enum.ExampleType
	Description string
	CreateTime  time.Time
	UpdateTime  time.Time
}

// ExampleQueryInterface example自定义查询接口
type ExampleQueryInterface interface {
	// SELECT * FROM @@table
	// WHERE
	// {{ if uuid != "" }} uuid = @uuid AND {{ end }}
	// {{ if name != "" }} name LIKE %@name% AND {{ end }}
	// 1=1
	Query(uuid string, name string) ([]*gen.T, error)
}

func (e *Example) TableName() string {
	return "example"
}

```

执行makefile指令
```shell
make gorm-gen
```

## 3. 实现你的grpc服务

参考`handler/api/v1/example.go`

```go
package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flamefatex/grpc-gateway-example/model"
	"github.com/flamefatex/grpc-gateway-example/model/query"
	"github.com/flamefatex/grpc-gateway-example/pkg/lib/errprx"
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
		err = errorx.ErrorfInternalServer("EXAMPLE_LIST_ERROR", "get example list failed, err: %s", err)
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

```

## 4. 注册grpc server和http endpoint

参考`handler/register_grpc.go`注册grpc server

```go
package handler

import (
	"context"

	v1 "github.com/flamefatex/grpc-gateway-example/handler/api/v1"
	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	"google.golang.org/grpc"
)

// ExecRegisterGrpcServiceServer 注册grpc服务处理器
func ExecRegisterGrpcServiceServer(ctx context.Context, grpcServer *grpc.Server) {
	// 注册
	proto_v1_example.RegisterExampleServiceServer(grpcServer, v1.NewExampleHandler())
}

```

参考`handler/register_grpc_gateway.go`注册http endpoint
```go
package handler

import (
	"context"

	proto_v1_example "github.com/flamefatex/grpc-gateway-example/proto/gen/go/api/v1/example"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GrpcGwRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

// ExecRegisterGrpcGatewayEndpoint 注册grpc gateway 端点
func ExecRegisterGrpcGatewayEndpoint(ctx context.Context) []GrpcGwRegister {
	regs := []GrpcGwRegister{
		// 注册http端点
		proto_v1_example.RegisterExampleServiceHandlerFromEndpoint,
	}
	return regs
}

```

## 5. 编译与访问

```shell
go run main.go
```

然后访问 `localhost:8082/api/v1/example/list`