// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
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
func (e *engine) CleanBuildExecutables(ctx context.Context) (int64, error) {
	logrus.Trace("clearing build executables in the database")

	switch e.config.Driver {
	case constants.DriverPostgres:
		res := e.client.Exec(CleanExecutablesPostgres)
		return res.RowsAffected, res.Error
	default:
		res := e.client.Exec(CleanExecutablesSqlite)
		return res.RowsAffected, res.Error
	}
}
