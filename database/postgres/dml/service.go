// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

const (
	// ListServices represents a query to
	// list all services in the database.
	ListServices = `
SELECT *
FROM services;
`

	// ListBuildServices represents a query to list
	// all services for a build_id in the database.
	ListBuildServices = `
SELECT *
FROM services
WHERE build_id = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// SelectBuildServicesCount represents a query to select
	// the count of services for a build_id in the database.
	SelectBuildServicesCount = `
SELECT count(*) as count
FROM services
WHERE build_id = ?
`

	// SelectServiceImagesCount represents a query to select
	// the count of an images appearances in the database.
	SelectServiceImagesCount = `
SELECT image, count(image) as count
FROM services
GROUP BY image
`

	// SelectServiceStatusesCount represents a query to select
	// the count of service status appearances in the database.
	SelectServiceStatusesCount = `
SELECT status, count(status) as count
FROM services
GROUP BY status;
`

	// SelectBuildService represents a query to select a
	// service for a build_id and number in the database.
	SelectBuildService = `
SELECT *
FROM services
WHERE build_id = ?
AND number = ?
LIMIT 1;
`

	// DeleteService represents a query to
	// remove a service from the database.
	DeleteService = `
DELETE
FROM services
WHERE id = ?;
`
)
