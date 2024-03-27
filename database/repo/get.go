// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
)

// GetRepo gets a repo by ID from the database.
func (e *engine) GetRepo(ctx context.Context, id int64) (*api.Repo, error) {
	e.logger.Tracef("getting repo %d from the database", id)

	// variable to store query results
	r := new(Repo)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableRepo).
		Preload("Owner").
		Where("id = ?", id).
		Take(r).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
	err = r.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted repos
		e.logger.Errorf("unable to decrypt repo %d: %v", id, err)

		// return the unencrypted repo
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.ToLibrary
		return r.ToAPI(), nil
	}

	// decrypt owner fields
	err = r.Owner.Decrypt(e.config.EncryptionKey)
	if err != nil {
		e.logger.Errorf("unable to decrypt repo owner %d: %v", id, err)

		// return the unencrypted repo
		return r.ToAPI(), nil
	}

	// return the decrypted repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.ToLibrary
	return r.ToAPI(), nil
}
