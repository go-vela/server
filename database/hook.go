// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// GetUser gets a user by unique ID from the database.
func (c *client) GetHook(number int, r *library.Repo) (*library.Hook, error) {
	logrus.Tracef("Getting webhook %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableHook).
		Raw(c.DML.HookService.Select["repo"], r.GetID(), number).
		Scan(h).Error

	return h.ToLibrary(), err
}

// GetLastHook gets the last hook by repo ID from the database.
func (c *client) GetLastHook(r *library.Repo) (*library.Hook, error) {
	logrus.Tracef("Getting last hook for repo %s from the database", r.GetFullName())

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableHook).
		Raw(c.DML.HookService.Select["last"], r.GetID()).
		Scan(h).Error

	// the record will not exist if it's a new repo
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	return h.ToLibrary(), err
}

// GetHookList gets a list of all webhooks from the database.
func (c *client) GetHookList() ([]*library.Hook, error) {
	logrus.Trace("Listing hooks from the database")

	// variable to store query results
	h := new([]database.Hook)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableHook).
		Raw(c.DML.HookService.List["all"]).
		Scan(h).Error

	// variable we want to return
	hooks := []*library.Hook{}
	// iterate through all query results
	for _, hook := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := hook

		// convert query result to library type
		hooks = append(hooks, tmp.ToLibrary())
	}

	return hooks, err
}

// GetRepoHookCount gets the count of webhooks by repo ID from the database.
func (c *client) GetRepoHookCount(r *library.Repo) (int64, error) {
	logrus.Tracef("Counting hooks for repo %s from the database", r.GetFullName())

	// variable to store query results
	var h []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableHook).
		Raw(c.DML.HookService.Select["count"], r.GetID()).
		Pluck("count", &h).Error

	return h[0], err
}

// GetRepoHookList gets a list of webhooks by repo ID from the database.
//
// nolint: dupl // ignore false positive
func (c *client) GetRepoHookList(r *library.Repo, page, perPage int) ([]*library.Hook, error) {
	logrus.Tracef("Listing hooks for repo %s from the database", r.GetFullName())

	// variable to store query results
	h := new([]database.Hook)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableHook).
		Raw(c.DML.HookService.List["repo"], r.GetID(), perPage, offset).
		Scan(h).Error

	// variable we want to return
	hooks := []*library.Hook{}
	// iterate through all query results
	for _, hook := range *h {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := hook

		// convert query result to library type
		hooks = append(hooks, tmp.ToLibrary())
	}

	return hooks, err
}

// CreateHook creates a new webhook in the database.
func (c *client) CreateHook(h *library.Hook) error {
	logrus.Tracef("Creating hook %s in the database", h.GetSourceID())

	// cast to database type
	hook := database.HookFromLibrary(h)

	// validate the necessary fields are populated
	err := hook.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableHook).
		Create(hook).Error
}

// UpdateHook updates a webhook in the database.
func (c *client) UpdateHook(h *library.Hook) error {
	logrus.Tracef("Updating hook %s in the database", h.GetSourceID())

	// cast to database type
	hook := database.HookFromLibrary(h)

	// validate the necessary fields are populated
	err := hook.Validate()
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
		Table(constants.TableHook).
		Where("id = ?", h.GetID()).
		Save(hook).Error
}

// DeleteHook deletes a webhook by unique ID from the database.
func (c *client) DeleteHook(id int64) error {
	logrus.Tracef("Deleting hook %d in the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableHook).
		Exec(c.DML.HookService.Delete, id).Error
}
