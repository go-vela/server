// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"errors"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// GetHook gets a hook by number and repo ID from the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) GetHook(number int, r *library.Repo) (*library.Hook, error) {
	logrus.Tracef("getting hook %s/%d from the database", r.GetFullName(), number)

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableHook).
		Raw(dml.SelectRepoHook, r.GetID(), number).
		Scan(h)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return h.ToLibrary(), result.Error
}

// GetLastHook gets the last hook by repo ID from the database.
func (c *client) GetLastHook(r *library.Repo) (*library.Hook, error) {
	logrus.Tracef("getting last hook for repo %s from the database", r.GetFullName())

	// variable to store query results
	h := new(database.Hook)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableHook).
		Raw(dml.SelectLastRepoHook, r.GetID()).
		Scan(h)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		// the record will not exist if it's a new repo
		return nil, nil
	}

	return h.ToLibrary(), result.Error
}

// CreateHook creates a new hook in the database.
func (c *client) CreateHook(h *library.Hook) error {
	logrus.Tracef("creating hook %d in the database", h.GetNumber())

	// cast to database type
	hook := database.HookFromLibrary(h)

	// validate the necessary fields are populated
	err := hook.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableHook).
		Create(hook).Error
}

// UpdateHook updates a hook in the database.
func (c *client) UpdateHook(h *library.Hook) error {
	logrus.Tracef("updating hook %d in the database", h.GetNumber())

	// cast to database type
	hook := database.HookFromLibrary(h)

	// validate the necessary fields are populated
	err := hook.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Sqlite.
		Table(constants.TableHook).
		Save(hook).Error
}

// DeleteHook deletes a hook by unique ID from the database.
func (c *client) DeleteHook(id int64) error {
	logrus.Tracef("deleting hook %d in the database", id)

	// send query to the database
	return c.Sqlite.
		Table(constants.TableHook).
		Exec(dml.DeleteHook, id).Error
}
