// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with update.go
package repo

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// CreateRepo creates a new repo in the database.
func (e *engine) CreateRepo(ctx context.Context, r *api.Repo) (*api.Repo, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("creating repo %s in the database", r.GetFullName())

	// cast the library type to database type
	repo := FromAPI(r)

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
	err = e.client.Table(constants.TableRepo).Create(repo).Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the repo
	err = repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// only log to preserve backwards compatibility
		e.logger.Errorf("unable to decrypt repo %d: %v", r.GetID(), err)

		return repo.ToAPI(), nil
	}

	// set owner to provided owner if creation successful
	result := repo.ToAPI()
	result.SetOwner(r.GetOwner())

	return result, nil
}
