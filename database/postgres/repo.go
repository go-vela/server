// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"errors"
	"fmt"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

// GetRepo gets a repo by org and name from the database.
func (c *client) GetRepo(org, name string) (*library.Repo, error) {
	logrus.Tracef("getting repo %s/%s from the database", org, name)

	// variable to store query results
	r := new(database.Repo)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableRepo).
		Raw(dml.SelectRepo, org, name).
		Scan(r)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// decrypt the fields for the repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
	err := r.Decrypt(c.config.EncryptionKey)
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted repos
		logrus.Errorf("unable to decrypt repo %s/%s: %v", org, name, err)

		// return the unencrypted repo
		return r.ToLibrary(), result.Error
	}

	// return the decrypted repo
	return r.ToLibrary(), result.Error
}

// CreateRepo creates a new repo in the database.
func (c *client) CreateRepo(r *library.Repo) error {
	logrus.Tracef("creating repo %s in the database", r.GetFullName())

	// cast to database type
	repo := database.RepoFromLibrary(r)

	// validate the necessary fields are populated
	err := repo.Validate()
	if err != nil {
		return err
	}

	// encrypt the fields for the repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Encrypt
	err = repo.Encrypt(c.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt repo %s: %v", r.GetFullName(), err)
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableRepo).
		Create(repo).Error
}

// UpdateRepo updates a repo in the database.
func (c *client) UpdateRepo(r *library.Repo) error {
	logrus.Tracef("updating repo %s in the database", r.GetFullName())

	// cast to database type
	repo := database.RepoFromLibrary(r)

	// validate the necessary fields are populated
	err := repo.Validate()
	if err != nil {
		return err
	}

	// encrypt the fields for the repo
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Encrypt
	err = repo.Encrypt(c.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt repo %s: %v", r.GetFullName(), err)
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableRepo).
		Save(repo).Error
}

// DeleteRepo deletes a repo by unique ID from the database.
func (c *client) DeleteRepo(id int64) error {
	logrus.Tracef("deleting repo %d in the database", id)

	// send query to the database
	return c.Postgres.
		Table(constants.TableRepo).
		Exec(dml.DeleteRepo, id).Error
}
