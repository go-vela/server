// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

const CleanExecutablesPostgres = `
	DELETE FROM build_executables 
	USING builds 
	WHERE builds.id = build_executables.build_id 
	AND builds.status = 'error';
`

const CleanExecutablesSqlite = `
	DELETE FROM build_executables
	WHERE build_id IN (
  		SELECT build_id FROM build_executables e
  		INNER JOIN builds b
    		ON e.build_id=b.id
  		WHERE b.status = 'error'
	);
`

// CleanBuildExecutables pops executables which have a corresponding build that was cleaned.
func (e *Engine) CleanBuildExecutables(ctx context.Context) (int64, error) {
	logrus.Trace("clearing build executables in the database")

	switch e.config.Driver {
	case constants.DriverPostgres:
		res := e.client.
			WithContext(ctx).
			Exec(CleanExecutablesPostgres)

		return res.RowsAffected, res.Error
	default:
		res := e.client.
			WithContext(ctx).
			Exec(CleanExecutablesSqlite)

		return res.RowsAffected, res.Error
	}
}
