// SPDX-License-Identifier: Apache-2.0

package database

// Close stops and terminates the connection to the database.
func (e *engine) Close() error {
	e.logger.Tracef("closing connection to the %s database", e.Driver())

	// capture database/sql database from gorm.io/gorm database
	_sql, err := e.client.DB()
	if err != nil {
		return err
	}

	return _sql.Close()
}
