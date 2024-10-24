// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"

	"gorm.io/gorm/clause"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// PopBuildExecutable pops a build executable by build_id from the database.
func (e *engine) PopBuildExecutable(ctx context.Context, id int64) (*api.BuildExecutable, error) {
	e.logger.Tracef("popping build executable for build %d", id)

	// variable to store query results
	b := new(types.BuildExecutable)

	// at the time of coding, GORM does not implement a version of Sqlite3 that supports RETURNING.
	// so we have to select and delete for the Sqlite driver.
	switch e.config.Driver {
	case constants.DriverPostgres:
		// send query to the database and store result in variable
		err := e.client.
			WithContext(ctx).
			Table(constants.TableBuildExecutable).
			Clauses(clause.Returning{}).
			Where("build_id = ?", id).
			Delete(b).
			Error
		if err != nil {
			return nil, err
		}

	case constants.DriverSqlite:
		// send query to the database and store result in variable
		err := e.client.
			WithContext(ctx).
			Table(constants.TableBuildExecutable).
			Where("id = ?", id).
			Take(b).
			Error
		if err != nil {
			return nil, err
		}

		// send query to the database to delete result just got
		err = e.client.
			WithContext(ctx).
			Table(constants.TableBuildExecutable).
			Delete(b).
			Error
		if err != nil {
			return nil, err
		}
	}

	// decrypt the fields for the build executable
	err := b.Decrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, err
	}

	// decompress data for the build executable
	err = b.Decompress()
	if err != nil {
		return nil, err
	}

	return b.ToAPI(), nil
}
