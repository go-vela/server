// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with create.go
package repo

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// UpdateRepo updates an existing repo in the database.
func (e *engine) UpdateRepo(r *library.Repo) (*library.Repo, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("creating repo %s in the database", r.GetFullName())

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#RepoFromLibrary
	repo := database.RepoFromLibrary(r)

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Validate
	err := repo.Validate()
	if err != nil {
		return nil, err
	}

	// encrypt the fields for the repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Encrypt
	err = repo.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("unable to encrypt repo %s: %w", r.GetFullName(), err)
	}

	// send query to the database
	err = e.client.Table(constants.TableRepo).Save(repo).Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the repo
	err = repo.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// only log to preserve backwards compatibility
		e.logger.Errorf("unable to decrypt repo %d: %v", r.GetID(), err)

		return repo.ToLibrary(), nil
	}

	return repo.ToLibrary(), nil
}
