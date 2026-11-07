package mysql

import (
	"context"

	"gorm.io/gorm"
)

// DB abstracts a GORM connection and lifecycle for easy mocking in tests.
type DB interface {
	GORM() *gorm.DB
	Ping(ctx context.Context) error
	Close() error
}

var _ DB = (*MysqlDB)(nil)

// GORM returns the underlying GORM handle (same as Gorm).
func (d *MysqlDB) GORM() *gorm.DB {
	return d.Gorm()
}
