// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetRepoList gets a list of all repos from the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) GetRepoList() ([]*library.Repo, error) {
	c.Logger.Trace("listing repos from the database")

	// variable to store query results
	r := new([]database.Repo)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableRepo).
		Raw(dml.ListRepos).
		Scan(r).Error
	if err != nil {
		return nil, err
	}

	// variable we want to return
	repos := []*library.Repo{}
	// iterate through all query results
	for _, repo := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// decrypt the fields for the repo
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
		err = tmp.Decrypt(c.config.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted repos
			c.Logger.Errorf("unable to decrypt repo %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		repos = append(repos, tmp.ToLibrary())
	}

	return repos, nil
}

// GetOrgRepoList gets a list of all repos by org from the database.
func (c *client) GetOrgRepoList(org string, filters map[string]string, page, perPage int) ([]*library.Repo, error) {
	c.Logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("listing repos for org %s from the database", org)

	// variable to store query results
	r := new([]database.Repo)

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableRepo).
		Where("org = ?", org).
		Where(filters).
		Limit(perPage).
		Offset(offset).
		Scan(r).Error
	if err != nil {
		return nil, err
	}

	// variable we want to return
	repos := []*library.Repo{}
	// iterate through all query results
	for _, repo := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// decrypt the fields for the repo
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
		err = tmp.Decrypt(c.config.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted repos
			c.Logger.Errorf("unable to decrypt repo %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		repos = append(repos, tmp.ToLibrary())
	}

	return repos, nil
}

// GetUserRepoList gets a list of all repos by user ID from the database.
func (c *client) GetUserRepoList(u *library.User, page, perPage int) ([]*library.Repo, error) {
	c.Logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("listing repos for user %s from the database", u.GetName())

	// variable to store query results
	r := new([]database.Repo)
	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableRepo).
		Raw(dml.ListUserRepos, u.GetID(), perPage, offset).
		Scan(r).Error
	if err != nil {
		return nil, err
	}

	// variable we want to return
	repos := []*library.Repo{}
	// iterate through all query results
	for _, repo := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// decrypt the fields for the repo
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Repo.Decrypt
		err = tmp.Decrypt(c.config.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted repos
			c.Logger.Errorf("unable to decrypt repo %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		repos = append(repos, tmp.ToLibrary())
	}

	return repos, nil
}
