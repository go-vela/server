// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"database/sql"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ListQueuedBuilds gets a list of all queued builds from the database.
func (c *client) ListQueuedBuilds() ([]*library.BuildQueue, error) {
	c.Logger.Trace("listing builds from the database")

	// variable to store query results
	b := new([]database.BuildQueue)

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableBuildQueue).
		Raw(dml.ListQueuedBuilds).
		Scan(b).Error

	// variable we want to return
	builds := []*library.BuildQueue{}
	// iterate through all query results
	for _, build := range *b {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := build

		// convert query result to library type
		builds = append(builds, tmp.ToLibrary())
	}

	return builds, err
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
	return c.Postgres.
		Table(constants.TableBuildQueue).
		Create(build).Error
}

func (c *client) PopQueuedBuild(tx *gorm.DB, id int64) error {
	c.Logger.WithFields(logrus.Fields{
		"item": id,
	}).Tracef("popping queued build %d in the database", id)

	var b library.BuildQueue

	// use transaction db if provided
	db := c.Postgres
	if tx != nil {
		db = tx
	}

	// todo: why doesnt raw query work?
	return db.
		Table(constants.TableBuildQueue).
		Where("build_id = ?", id).
		Delete(&b).
		Raw(dml.DeleteQueuedBuild, id).Error
}

func (c *client) Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return c.Postgres.Transaction(fc, opts...)
}

func (c *client) TryTransactionLock(txID int, tx *gorm.DB) error {
	// use transaction db if provided
	db := c.Postgres
	if tx != nil {
		db = tx
	}

	err := db.Exec("SELECT PG_TRY_ADVISORY_XACT_LOCK(789);").Error
	if err != nil {
		return err
	}

	return nil
}