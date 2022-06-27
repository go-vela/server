// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetRepo gets a repo by ID from the database.
func (e *engine) GetRepo(id int64) (*library.Repo, error) {
	e.logger.Tracef("getting repo %d from the database", id)

	// variable to store query results
	r := new(database.Repo)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableRepo).
		Where("id = ?", id).
		Limit(1).
		Scan(r).
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
		return r.ToLibrary(), nil
	}

	// return the decrypted repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.ToLibrary
	return r.ToLibrary(), nil
}
