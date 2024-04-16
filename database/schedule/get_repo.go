// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// GetScheduleForRepo gets a schedule by repo ID and name from the database.
func (e *engine) GetScheduleForRepo(ctx context.Context, r *api.Repo, name string) (*api.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"org":      r.GetOrg(),
		"repo":     r.GetName(),
		"schedule": name,
	}).Tracef("getting schedule %s/%s from the database", r.GetFullName(), name)

	// variable to store query results
	s := new(Schedule)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSchedule).
		Where("repo_id = ?", r.GetID()).
		Where("name = ?", name).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	return s.ToAPI(), nil
}
