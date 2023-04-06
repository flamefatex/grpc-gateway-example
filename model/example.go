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
