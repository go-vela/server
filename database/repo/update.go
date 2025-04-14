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
func (e *Engine) UpdateRepo(ctx context.Context, r *api.Repo) (*api.Repo, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("creating repo %s", r.GetFullName())

	// cast the API type to database type
	repo := types.RepoFromAPI(r)

	// validate the necessary fields are populated
	err := repo.Validate()
	if err != nil {
		return nil, err
	}

	// encrypt the fields for the repo
	err = repo.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt repo %s: %w", r.GetFullName(), err)
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Save(repo).Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the repo
	err = repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// only log to preserve backwards compatibility
		e.logger.Errorf("unable to decrypt repo %d: %v", r.GetID(), err)
	}

	// set owner to provided owner if creation successful
	result := repo.ToAPI()
	result.SetOwner(r.GetOwner())

	return result, nil
}
