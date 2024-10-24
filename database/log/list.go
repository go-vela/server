// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListLogs gets a list of all logs from the database.
func (e *engine) ListLogs(ctx context.Context) ([]*api.Log, error) {
	e.logger.Trace("listing all logs")

	// variables to store query results and return value
	count := int64(0)
	l := new([]types.Log)
	logs := []*api.Log{}

	// count the results
	count, err := e.CountLogs(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return logs, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Find(&l).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, log := range *l {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := log

		// decompress log data
		err = tmp.Decompress()
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch uncompressed logs
			e.logger.Errorf("unable to decompress logs: %v", err)
		}

		// convert query result to API type
		logs = append(logs, tmp.ToAPI())
	}

	return logs, nil
}
