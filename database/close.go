// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

// Close stops and terminates the connection to the database.
func (e *engine) Close() error {
	e.Logger.Tracef("closing connection to the %s database", e.Driver())

	// capture database/sql database from gorm.io/gorm database
	_sql, err := e.Database.DB()
	if err != nil {
		return err
	}

	return _sql.Close()
}
