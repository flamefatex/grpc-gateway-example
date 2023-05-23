package model

import (
	"time"

	proto_enum "github.com/flamefatex/grpc-gateway-example/proto/gen/go/enumeration"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type Example struct {
	Id          string
	Name        string
	Type        proto_enum.ExampleType
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

// ExampleQueryInterface example自定义查询接口
type ExampleQueryInterface interface {
	// SELECT * FROM @@table
	// WHERE
	// {{ if id != "" }} id = @id AND {{ end }}
	// {{ if name != "" }} name LIKE %@name% AND {{ end }}
	// 1=1
	Query(id string, name string) ([]*gen.T, error)
}

func (m *Example) TableName() string {
	return "example"
}
