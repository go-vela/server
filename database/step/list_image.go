// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"database/sql"

	"github.com/go-vela/server/constants"
)

// ListStepImageCount gets a list of all step images and the count of their occurrence from the database.
func (e *Engine) ListStepImageCount(ctx context.Context) (map[string]float64, error) {
	e.logger.Tracef("getting count of all images for steps")

	// variables to store query results and return value
	s := []struct {
		Image sql.NullString
		Count sql.NullInt32
	}{}
	images := make(map[string]float64)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableStep).
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
