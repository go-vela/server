// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateRepo updates an existing repo in the database.
func (e *Engine) PartialUpdateRepo(ctx context.Context, r *api.Repo) error {
	e.logger.WithFields(logrus.Fields{
		"id": r.GetID(),
	}).Tracef("updating repo %d", r.GetID())

	if r.GetID() == 0 {
		return fmt.Errorf("repo ID must be set")
	}

	repo := types.RepoFromAPI(r)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Updates(repo).Error
}
