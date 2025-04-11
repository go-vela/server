// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListLogs gets a list of all logs from the database.
func (e *Engine) ListLogs(ctx context.Context) ([]*api.Log, error) {
	e.logger.Trace("listing all logs")

	// variables to store query results and return value
	l := new([]types.Log)
	logs := []*api.Log{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Find(&l).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, log := range *l {
		// decompress log data
		err = log.Decompress()
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch uncompressed logs
			e.logger.Errorf("unable to decompress logs: %v", err)
		}

		// convert query result to API type
		logs = append(logs, log.ToAPI())
	}

	return logs, nil
}
