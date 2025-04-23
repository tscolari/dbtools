package gormtest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tscolari/servicetools/database/dbtest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is a utility wrapper around dbtest.DB to return a gorm.DB object.
func DB(t *testing.T, migrationsPath, name string) (*gorm.DB, func()) {
	db, closer := dbtest.DB(t, migrationsPath, name)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	require.NoError(t, err)

	return gormDB, closer
}
