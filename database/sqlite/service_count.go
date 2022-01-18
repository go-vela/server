// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetBuildServiceCount gets a count of all services by build ID from the database.
func (c *client) GetBuildServiceCount(b *library.Build) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting count of services for build %d from the database", b.GetNumber())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.SelectBuildServicesCount, b.GetID()).
		Pluck("count", &s).Error

	return s, err
}

// GetServiceImageCount gets a count of all service images
// and the count of their occurrence in the database.
func (c *client) GetServiceImageCount() (map[string]float64, error) {
	c.Logger.Tracef("getting count of all images for services from the database")

	type imageCount struct {
		Image string
		Count int
	}

	// variable to store query results
	images := new([]imageCount)
	counts := make(map[string]float64)

	// send query to the database and store result in variable
	err := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.SelectServiceImagesCount).
		Scan(images).Error

	for _, image := range *images {
		counts[image.Image] = float64(image.Count)
	}

	return counts, err
}

// GetServiceStatusCount gets a list of all service statuses
// and the count of their occurrence in the database.
func (c *client) GetServiceStatusCount() (map[string]float64, error) {
	c.Logger.Trace("getting count of all statuses for services from the database")

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
	err := c.Sqlite.
		Table(constants.TableService).
		Raw(dml.SelectServiceStatusesCount).
		Scan(s).Error

	for _, status := range *s {
		counts[status.Status] = float64(status.Count)
	}

	return counts, err
}
