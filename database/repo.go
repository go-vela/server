// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetRepo gets a repo by org and name from the database.
func (c *client) GetRepo(org, name string) (*library.Repo, error) {
	logrus.Tracef("Getting repo %s/%s from the database", org, name)

	// variable to store query results
	r := new(database.Repo)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableRepo).
		Raw(c.DML.RepoService.Select["repo"], org, name).
		Scan(r).Error

	return r.ToLibrary(), err
}

// GetRepoList gets a list of all repos from the database.
func (c *client) GetRepoList() ([]*library.Repo, error) {
	logrus.Trace("Listing repos from the database")

	// variable to store query results
	r := new([]database.Repo)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableRepo).
		Raw(c.DML.RepoService.List["all"]).
		Scan(r).Error

	// variable we want to return
	repos := []*library.Repo{}
	// iterate through all query results
	for _, repo := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// convert query result to library type
		repos = append(repos, tmp.ToLibrary())
	}

	return repos, err
}

// GetRepoCount gets a count of all repos from the database.
func (c *client) GetRepoCount() (int64, error) {
	logrus.Trace("Counting repos in the database")

	// variable to store query results
	var r []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableRepo).
		Raw(c.DML.RepoService.Select["count"]).
		Pluck("count", &r).Error

	return r[0], err
}

// GetUserRepoList gets a list of all repos by user ID from the database.
func (c *client) GetUserRepoList(u *library.User, page, perPage int) ([]*library.Repo, error) {
	logrus.Tracef("Listing repos for user %s from the database", u.GetName())

	// variable to store query results
	r := new([]database.Repo)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableRepo).
		Raw(c.DML.RepoService.List["user"], u.GetID(), perPage, offset).
		Scan(r).Error

	// variable we want to return
	repos := []*library.Repo{}
	// iterate through all query results
	for _, repo := range *r {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := repo

		// convert query result to library type
		repos = append(repos, tmp.ToLibrary())
	}

	return repos, err
}

// GetUserRepoCount gets a count of all repos for a specific user from the database.
func (c *client) GetUserRepoCount(u *library.User) (int64, error) {
	logrus.Tracef("Counting repos for user %s in the database", u.GetName())

	// variable to store query results
	var r []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableRepo).
		Raw(c.DML.RepoService.Select["countByUser"], u.GetID()).
		Pluck("count", &r).Error

	return r[0], err
}

// CreateRepo creates a new repo in the database.
func (c *client) CreateRepo(r *library.Repo) error {
	logrus.Tracef("Creating repo %s in the database", r.GetFullName())

	// cast to database type
	repo := database.RepoFromLibrary(r)

	// validate the necessary fields are populated
	err := repo.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableRepo).
		Create(repo).Error
}

// UpdateRepo updates a repo in the database.
func (c *client) UpdateRepo(r *library.Repo) error {
	logrus.Tracef("Updating repo %s in the database", r.GetFullName())

	// cast to database type
	repo := database.RepoFromLibrary(r)

	// validate the necessary fields are populated
	err := repo.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	//
	// .Update(repo) doesnt allow setting
	// select booleans to false.
	//
	// ref: https://github.com/jinzhu/gorm/issues/202
	return c.Database.
		Table(constants.TableRepo).
		Where("id = ?", r.GetID()).
		Save(repo).Error
}

// DeleteRepo deletes a repo by unique ID from the database.
func (c *client) DeleteRepo(id int64) error {
	logrus.Tracef("Deleting repo %d in the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableRepo).
		Raw(c.DML.RepoService.Delete, id).Error
}
