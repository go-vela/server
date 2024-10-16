// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"

	"gorm.io/gorm/clause"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// PopBuildExecutable pops a build executable by build_id from the database.
func (e *engine) PopBuildExecutable(ctx context.Context, id int64) (*library.BuildExecutable, error) {
	e.logger.Tracef("popping build executable for build %d", id)

	// variable to store query results
	b := new(database.BuildExecutable)

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
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
	err := b.Decrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, err
	}

	// decompress data for the build executable
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildExecutable.Decompress
	err = b.Decompress()
	if err != nil {
		return nil, err
	}

	// return the decompressed build executable
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#BuildExecutable.ToLibrary
	return b.ToLibrary(), nil
}
