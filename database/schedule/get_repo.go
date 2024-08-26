// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetScheduleForRepo gets a schedule by repo ID and name from the database.
func (e *engine) GetScheduleForRepo(ctx context.Context, r *api.Repo, name string) (*api.Schedule, error) {
	e.logger.WithFields(logrus.Fields{
		"org":      r.GetOrg(),
		"repo":     r.GetName(),
		"schedule": name,
	}).Tracef("getting schedule %s/%s", r.GetFullName(), name)

	// variable to store query results
	s := new(types.Schedule)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Preload("Repo").
		Preload("Repo.Owner").
		Where("repo_id = ?", r.GetID()).
		Where("name = ?", name).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt hash value for repo
	err = s.Repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo %d: %v", s.Repo.ID.Int64, err)
	}

	return s.ToAPI(), nil
}
