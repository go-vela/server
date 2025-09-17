// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// testPostgresGormInit initializes a Gorm DB for postgres testing.
//
// TODO: remove this when Gorm updates a bug they introduced in 1.30.1 where the
// special dialector config is overwritten in the Open(..) function.
func TestPostgresGormInit(sql *sql.DB) (*gorm.DB, error) {
	_postgres, err := gorm.Open(postgres.New(postgres.Config{Conn: sql}))
	if err != nil {
		return nil, err
	}

	_postgres = _postgres.Session(&gorm.Session{SkipDefaultTransaction: true})

	return _postgres, nil
}
