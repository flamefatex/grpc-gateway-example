package gorm_gen

import (
	"github.com/flamefatex/grpc-gateway-example/model"
	"gorm.io/gen"
)

// Apply 应用sql接口
func Apply(g *gen.Generator) {
	g.ApplyBasic(
		model.Example{},
	)

	// Id
	g.ApplyInterface(func(IdQueryInterface) {},
		model.Example{},
	)

	// Example
	g.ApplyInterface(func(queryInterface model.ExampleQueryInterface) {}, model.Example{})
}
