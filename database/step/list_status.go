// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"database/sql"

	"github.com/go-vela/server/constants"
)

// ListStepStatusCount gets a list of all step statuses and the count of their occurrence from the database.
func (e *Engine) ListStepStatusCount(ctx context.Context) (map[string]float64, error) {
	e.logger.Tracef("getting count of all statuses for steps")

	// variables to store query results and return value
	s := []struct {
		Status sql.NullString
		Count  sql.NullInt32
	}{}
	statuses := map[string]float64{
		"pending": 0,
		"failure": 0,
		"killed":  0,
		"running": 0,
		"success": 0,
	}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
		Select("status", " count(status) as count").
		Group("status").
		Find(&s).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, value := range s {
		// check if the status returned is not empty
		if value.Status.Valid {
			statuses[value.Status.String] = float64(value.Count.Int32)
		}
	}

	return statuses, nil
}
