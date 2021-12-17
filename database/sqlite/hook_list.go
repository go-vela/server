// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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

// GetHookList gets a list of all hooks from the database.
func (c *client) GetHookList() ([]*library.Hook, error) {
	c.Logger.Trace("listing hooks from the database")

	// variable to store query results
	h := new([]database.Hook)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableHook).
		Raw(dml.ListHooks).
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

// GetRepoHookList gets a list of hooks by repo ID from the database.
func (c *client) GetRepoHookList(r *library.Repo, page, perPage int) ([]*library.Hook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("listing hooks for repo %s from the database", r.GetFullName())

	// variable to store query results
	h := new([]database.Hook)
	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableHook).
		Raw(dml.ListRepoHooks, r.GetID(), perPage, offset).
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
