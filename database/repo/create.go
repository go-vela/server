// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"fmt"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CreateRepo creates a new repo in the database.
func (e *engine) CreateRepo(r *library.Repo) error {
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
		return err
	}

	// encrypt the fields for the repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Encrypt
	err = repo.Encrypt(e.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt repo %s: %w", r.GetFullName(), err)
	}

	// send query to the database
	return e.client.
		Table(constants.TableRepo).
		Create(repo).
		Error
}
