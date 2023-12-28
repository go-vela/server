// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetRepoForOrg gets a repo by org and repo name from the database.
func (e *engine) GetRepoForOrg(ctx context.Context, org, name string) (*library.Repo, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"repo": name,
	}).Tracef("getting repo %s/%s from the database", org, name)

	// variable to store query results
	r := new(database.Repo)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableRepo).
		Where("org = ?", org).
		Where("name = ?", name).
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
		e.logger.Errorf("unable to decrypt repo %s/%s: %v", org, name, err)

		// return the unencrypted repo
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.ToLibrary
		return r.ToLibrary(), nil
	}

	// return the decrypted repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.ToLibrary
	return r.ToLibrary(), nil
}
