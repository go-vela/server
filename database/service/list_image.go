// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"database/sql"

	"github.com/go-vela/types/constants"
)

// ListServiceImageCount gets a list of all service images and the count of their occurrence from the database.
func (e *engine) ListServiceImageCount() (map[string]float64, error) {
	e.logger.Tracef("getting count of all images for services from the database")

	// variables to store query results and return value
	s := []struct {
		Image sql.NullString
		Count sql.NullInt32
	}{}
	images := make(map[string]float64)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableService).
		Select("image", " count(image) as count").
		Group("image").
		Find(&s).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, value := range s {
		// check if the image returned is not empty
		if value.Image.Valid {
			images[value.Image.String] = float64(value.Count.Int32)
		}
	}

	return images, nil
}
