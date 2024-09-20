// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetRepoForOrg gets a repo by org and repo name from the database.
func (e *engine) GetRepoForOrg(ctx context.Context, fullName string) (*api.Repo, error) {
	// e.logger.WithFields(logrus.Fields{
	// 	"org":  org,
	// 	"repo": name,
	// }).Tracef("getting repo %s/%s", org, name)

	// variable to store query results
	r := new(types.Repo)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableRepo).
		Preload("Owner").
		Where("full_name = ?", fullName).
		Take(r).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the repo
	err = r.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted repos
		e.logger.Errorf("unable to decrypt repo %s: %v", fullName, err)

		// return the unencrypted repo
		return r.ToAPI(), nil
	}

	// return the decrypted repo
	return r.ToAPI(), nil
}
