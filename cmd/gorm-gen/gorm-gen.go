package main

import (
	"flag"

	gorm_gen "github.com/flamefatex/grpc-gateway-example/model/gorm-gen"
	"gorm.io/gen"
)

var (
	path = flag.String("path", "model/query", "The output path of query files.")
)

type IDAppendMethod interface {
	// SELECT * FROM @@table WHERE id = @id
	GetById(id int64) (*gen.T, error)

	// DELETE FROM @@table WHERE id = @id
	DeleteById(id int64) (gen.RowsAffected, error)
}

type UUIDAppendMethod interface {
	// SELECT * FROM @@table WHERE uuid = @uuid
	GetByUuid(uuid string) (*gen.T, error)

	// DELETE FROM @@table WHERE uuid = @uuid
	DeleteByUuid(uuid string) (gen.RowsAffected, error)
}

func main() {
	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath: *path,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// Generate basic type-safe DAO API
	gorm_gen.Apply(g)

	g.Execute()
}
