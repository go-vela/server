// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

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
WHERE build_id = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// SelectBuildStepsCount represents a query to select
	// the count of steps for a build_id in the database.
	SelectBuildStepsCount = `
SELECT count(*) as count
FROM steps
WHERE build_id = ?
`

	// SelectStepImagesCount represents a query to select
	// the count of an images appearances in the database.
	SelectStepImagesCount = `
SELECT image, count(image) as count
FROM steps
GROUP BY image;
`

	// SelectStepStatusesCount represents a query to select
	// the count of step status' appearances in the database.
	SelectStepStatusesCount = `
SELECT status, count(status) as count
FROM steps
GROUP BY status;
`

	// SelectBuildStep represents a query to select a
	// step for a build_id and number in the database.
	SelectBuildStep = `
SELECT *
FROM steps
WHERE build_id = ?
AND number = ?
LIMIT 1;
`

	// DeleteStep represents a query to
	// remove a step from the database.
	DeleteStep = `
DELETE
FROM steps
WHERE id = ?;
`
)
