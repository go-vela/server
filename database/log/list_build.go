// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListLogsForBuild gets a list of logs by build ID from the database.
func (e *Engine) ListLogsForBuild(ctx context.Context, b *api.Build, page, perPage int) ([]*api.Log, error) {
	e.logger.Tracef("listing logs for build %d", b.GetID())

	// variables to store query results and return value
	l := new([]types.Log)
	logs := []*api.Log{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableLog).
		Where("build_id = ?", b.GetID()).
		Order("service_id ASC NULLS LAST").
		Order("step_id ASC").
		Limit(perPage).
		Offset(offset).
		Find(&l).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, log := range *l {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := log

		// decompress log data for the build
		err = tmp.Decompress()
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch uncompressed logs
			e.logger.Errorf("unable to decompress logs for build %d: %v", b.GetID(), err)
		}

		// convert query result to API type
		logs = append(logs, tmp.ToAPI())
	}

	return logs, nil
}
