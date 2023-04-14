package gorm_gen

import "gorm.io/gen"

type IdQueryInterface interface {
	// SELECT * FROM @@table WHERE id = @id
	GetById(id string) (*gen.T, error)

	// DELETE FROM @@table WHERE id = @id
	DeleteById(id string) (gen.RowsAffected, error)
}
