// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// CleanBuilds updates builds to an error with a provided message with a created timestamp prior to a defined moment.
func (e *engine) CleanBuilds(ctx context.Context, msg string, statuses []string, before int64) ([]*api.Build, int64, error) {
	logrus.Tracef("cleaning pending or running builds in the database created prior to %d", before)

	b := new(api.Build)
	b.SetStatus(constants.StatusError)
	b.SetError(msg)
	b.SetFinished(time.Now().UTC().Unix())

	build := types.BuildFromAPI(b)
	builds := new([]types.Build)

	statusQuery := e.client.Table(constants.TableBuild)
	for _, status := range statuses {
		statusQuery = statusQuery.Or("status = ?", status)
	}

	var (
		count     int64
		apiBuilds []*api.Build
	)

	switch e.config.Driver {
	case constants.DriverPostgres:
		result := e.client.
			Table(constants.TableBuild).
			Model(&builds).
			Clauses(clause.Returning{}).
			Where("created < ?", before).
			Where(statusQuery).
			Updates(build)
		if result.Error != nil {
			return nil, 0, result.Error
		}

		count = result.RowsAffected
	case constants.DriverSqlite:
		result := e.client.
			Table(constants.TableBuild).
			Where("created < ?", before).
			Where(statusQuery).
			Find(&builds)

		if result.Error != nil {
			return nil, 0, result.Error
		}

		result = e.client.
			Table(constants.TableBuild).
			Where("created < ?", before).
			Where(statusQuery).
			Updates(build)
		if result.Error != nil {
			return nil, 0, result.Error
		}

		count = result.RowsAffected
	}

	// iterate through all query results
	for _, bld := range *builds {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := bld

		apiBuilds = append(apiBuilds, tmp.ToAPI())
	}

	return apiBuilds, count, nil
}
