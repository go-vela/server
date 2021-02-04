// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetStep gets a step by number and build ID from the database.
func (c *client) GetStep(number int, b *library.Build) (*library.Step, error) {
	logrus.Tracef("Getting step %d for build %d from the database", number, b.GetNumber())

	// variable to store query results
	s := new(database.Step)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableStep).
		Raw(c.DML.StepService.Select["build"], b.GetID(), number).
		Scan(s).Error

	return s.ToLibrary(), err
}

// GetStepList gets a list of all steps from the database.
func (c *client) GetStepList() ([]*library.Step, error) {
	logrus.Trace("Listing steps from the database")

	// variable to store query results
	s := new([]database.Step)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableStep).
		Raw(c.DML.StepService.List["all"]).
		Scan(s).Error

	// variable we want to return
	steps := []*library.Step{}
	// iterate through all query results
	for _, step := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := step

		// convert query result to library type
		steps = append(steps, tmp.ToLibrary())
	}

	return steps, err
}

// GetBuildStepList gets a list of all steps by build ID from the database.
//
// nolint: dupl // ignore false positive
func (c *client) GetBuildStepList(b *library.Build, page, perPage int) ([]*library.Step, error) {
	logrus.Tracef("Listing steps for build %d from the database", b.GetNumber())

	// variable to store query results
	s := new([]database.Step)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableStep).
		Raw(c.DML.StepService.List["build"], b.GetID(), perPage, offset).
		Scan(s).Error

	// variable we want to return
	steps := []*library.Step{}
	// iterate through all query results
	for _, step := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := step

		// convert query result to library type
		steps = append(steps, tmp.ToLibrary())
	}

	return steps, err
}

// GetBuildStepCount gets a count of all steps by build ID from the database.
func (c *client) GetBuildStepCount(b *library.Build) (int64, error) {
	logrus.Tracef("Counting build steps for build %d in the database", b.GetNumber())

	// variable to store query results
	var r []int64

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableStep).
		Raw(c.DML.StepService.Select["count"], b.GetID()).
		Pluck("count", &r).Error

	return r[0], err
}

// GetStepImageCount gets a list of all step images
// and the count of their occurrence in the database.
func (c *client) GetStepImageCount() (map[string]float64, error) {
	logrus.Tracef("Counting images for steps in the database")

	type imageCount struct {
		Image string
		Count int
	}

	// variable to store query results
	images := new([]imageCount)
	counts := make(map[string]float64)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableStep).
		Raw(c.DML.StepService.Select["count-images"]).
		Scan(images).Error

	for _, image := range *images {
		counts[image.Image] = float64(image.Count)
	}

	return counts, err
}

// GetStepStatusCount gets a list of all step statuses
// and the count of their occurrence in the database.
func (c *client) GetStepStatusCount() (map[string]float64, error) {
	logrus.Trace("Counting the total of each status for steps in the database")

	type statusCount struct {
		Status string
		Count  int
	}

	// variable to store query results
	s := new([]statusCount)
	counts := map[string]float64{
		"pending": 0,
		"failure": 0,
		"killed":  0,
		"running": 0,
		"success": 0,
	}

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableStep).
		Raw(c.DML.StepService.Select["count-statuses"]).
		Scan(s).Error

	for _, status := range *s {
		counts[status.Status] = float64(status.Count)
	}

	return counts, err
}

// CreateStep creates a new step in the database.
func (c *client) CreateStep(s *library.Step) error {
	logrus.Tracef("Creating step %s in the database", s.GetName())

	// cast to database type
	step := database.StepFromLibrary(s)

	// validate the necessary fields are populated
	err := step.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableStep).
		Create(step).Error
}

// UpdateStep updates a step in the database.
func (c *client) UpdateStep(s *library.Step) error {
	logrus.Tracef("Updating step %s in the database", s.GetName())

	// cast to database type
	step := database.StepFromLibrary(s)

	// validate the necessary fields are populated
	err := step.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return c.Database.
		Table(constants.TableStep).
		Where("id = ?", s.GetID()).
		Update(step).Error
}

// DeleteStep deletes a step by unique ID from the database.
func (c *client) DeleteStep(id int64) error {
	logrus.Tracef("Deleting step %d from the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableStep).
		Exec(c.DML.StepService.Delete, id).Error
}
