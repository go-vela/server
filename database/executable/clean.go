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
func (e *engine) CleanBuildExecutables(ctx context.Context) error {
	logrus.Trace("clearing build executables in the database")

	switch e.config.Driver {
	case constants.DriverPostgres:
		return e.client.Exec(CleanExecutablesPostgres).Error

	default:
		return e.client.Exec(CleanExecutablesSqlite).Error
	}
}
