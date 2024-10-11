// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// LastHookForRepo gets the last hook by repo ID from the database.
func (e *engine) LastHookForRepo(ctx context.Context, r *api.Repo) (*api.Hook, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("getting last hook for repo %s", r.GetFullName())

	// variable to store query results
	h := new(types.Hook)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableHook).
		Preload("Build").
		Where("repo_id = ?", r.GetID()).
		Order("number DESC").
		Take(h).
		Error
	if err != nil {
		// check if the query returned a record not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// the record will not exist if it is a new repo
			return nil, nil
		}

		return nil, err
	}

	result := h.ToAPI()
	result.SetRepo(r)

	return result, nil
}
