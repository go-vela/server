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

func (c *client) ListQueuedBuilds() ([]*library.BuildQueue, error) {
	c.Logger.Trace("getting pending and running builds from the database")

	// variable to store query results
	b := new([]database.BuildQueue)

	// send query to the database and store result in variable
	result := c.Sqlite.
		Table(constants.TableBuildQueue).
		Raw(dml.ListQueuedBuilds).
		Scan(b)

	// variable we want to return
	builds := []*library.BuildQueue{}

	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, result.Error
}

func (c *client) CreateQueuedBuild(b *library.BuildQueue) error {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("creating queued build %d in the database", b.GetBuildID())

	// cast to database type
	build := database.BuildQueueFromLibrary(b)

	// validate the necessary fields are populated
	// err := build.Validate()
	// if err != nil {
	// 	return err
	// }

	// send query to the database
	return c.Sqlite.
		Table(constants.TableBuildQueue).
		Create(build).Error
}

func (c *client) PopQueuedBuild(id int64) error {
	c.Logger.WithFields(logrus.Fields{
		"item": id,
	}).Tracef("popping queued build %d in the database", id)

	return c.Sqlite.
		Table(constants.TableBuildQueue).
		Delete(id).Error
}
