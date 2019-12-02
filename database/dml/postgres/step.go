// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

const (
	// ListSteps represents a query to
	// list all steps in the database.
	ListSteps = `
SELECT *
FROM steps;
`

	// ListBuildSteps represents a query to list
	// all steps for a build_id in the database.
	ListBuildSteps = `
SELECT *
FROM steps
WHERE build_id = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;
`

	// SelectBuildStepsCount represents a query to select
	// the count of steps for a build_id in the database.
	SelectBuildStepsCount = `
SELECT count(*) as count
FROM steps
WHERE build_id = $1
`

	// SelectStepImagesCount represents a query to select
	// the count of an images appearances in the database.
	SelectStepImagesCount = `
SELECT image, COUNT(image)
FROM steps
GROUP BY image
`

	// SelectBuildStep represents a query to select a
	// step for a build_id and number in the database.
	SelectBuildStep = `
SELECT *
FROM steps
WHERE build_id = $1
AND number = $2
LIMIT 1;
`

	// DeleteStep represents a query to
	// remove a step from the database.
	DeleteStep = `
DELETE
FROM steps
WHERE id = $1;
`
)

// createStepService is a helper function to return
// a service for interacting with the steps table.
func createStepService() *Service {
	return &Service{
		List: map[string]string{
			"all":   ListSteps,
			"build": ListBuildSteps,
		},
		Select: map[string]string{
			"build":       SelectBuildStep,
			"count":       SelectBuildStepsCount,
			"image-count": SelectStepImagesCount,
		},
		Delete: DeleteStep,
	}
}
