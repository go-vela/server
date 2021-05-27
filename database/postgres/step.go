// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"errors"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

// GetStep gets a step by number and build ID from the database.
func (c *client) GetStep(number int, b *library.Build) (*library.Step, error) {
	logrus.Tracef("getting step %d for build %d from the database", number, b.GetNumber())

	// variable to store query results
	s := new(database.Step)

	// send query to the database and store result in variable
	result := c.Postgres.
		Table(constants.TableStep).
		Raw(dml.SelectBuildStep, b.GetID(), number).
		Scan(s)

	// check if the query returned a record not found error or no rows were returned
	if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return s.ToLibrary(), result.Error
}

// CreateStep creates a new step in the database.
func (c *client) CreateStep(s *library.Step) error {
	logrus.Tracef("creating step %s in the database", s.GetName())

	// cast to database type
	step := database.StepFromLibrary(s)

	// validate the necessary fields are populated
	err := step.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableStep).
		Create(step).Error
}

// UpdateStep updates a step in the database.
func (c *client) UpdateStep(s *library.Step) error {
	logrus.Tracef("updating step %s in the database", s.GetName())

	// cast to database type
	step := database.StepFromLibrary(s)

	// validate the necessary fields are populated
	err := step.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableStep).
		Save(step).Error
}

// DeleteStep deletes a step by unique ID from the database.
func (c *client) DeleteStep(id int64) error {
	logrus.Tracef("deleting step %d from the database", id)

	// send query to the database
	return c.Postgres.
		Table(constants.TableStep).
		Exec(dml.DeleteStep, id).Error
}
