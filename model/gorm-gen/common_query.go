package gorm_gen

import "gorm.io/gen"

type IdQueryInterface interface {
	// SELECT * FROM @@table WHERE id = @id
	GetById(id int64) (*gen.T, error)

	// DELETE FROM @@table WHERE id = @id
	DeleteById(id int64) (gen.RowsAffected, error)
}

type UuidQueryInterface interface {
	// SELECT * FROM @@table WHERE uuid = @uuid
	GetByUuid(uuid string) (*gen.T, error)

	// DELETE FROM @@table WHERE uuid = @uuid
	DeleteByUuid(uuid string) (gen.RowsAffected, error)
}
