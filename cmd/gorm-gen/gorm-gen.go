package main

import (
	"flag"

	gorm_gen "github.com/flamefatex/grpc-gateway-example/model/gorm-gen"
	"gorm.io/gen"
)

var (
	path = flag.String("path", "model/query", "The output path of query files.")
)

func main() {
	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath: *path,
		Mode:    gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// Generate basic type-safe DAO API
	gorm_gen.Apply(g)

	g.Execute()
}
