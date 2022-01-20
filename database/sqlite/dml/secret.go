// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package dml

const (
	// ListSecrets represents a query to
	// list all secrets in the database.
	//
	// nolint: gosec // ignore false positive
	ListSecrets = `
SELECT *
FROM secrets;
`

	// ListOrgSecrets represents a query to list all
	// secrets for a type and org in the database.
	//
	// nolint: gosec // ignore false positive
	ListOrgSecrets = `
SELECT *
FROM secrets
WHERE type = 'org'
AND org = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// ListRepoSecrets represents a query to list all
	// secrets for a type, org and repo in the database.
	//
	// nolint: gosec // ignore false positive
	ListRepoSecrets = `
SELECT *
FROM secrets
WHERE type = 'repo'
AND org = ?
AND repo = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// ListSharedSecrets represents a query to list all
	// secrets for a type, org and team in the database.
	//
	// nolint: gosec // ignore false positive
	ListSharedSecrets = `
SELECT *
FROM secrets
WHERE type = 'shared'
AND org = ?
AND team = ?
ORDER BY id DESC
LIMIT ?
OFFSET ?;
`

	// SelectOrgSecretsCount represents a query to select the
	// count of org secrets for an org in the database.
	//
	// nolint: gosec // ignore false positive
	SelectOrgSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'org'
AND org = ?;
`

	// SelectRepoSecretsCount represents a query to select the
	// count of repo secrets for an org and repo in the database.
	//
	// nolint: gosec // ignore false positive
	SelectRepoSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'repo'
AND org = ?
AND repo = ?;
`

	// SelectSharedSecretsCount represents a query to select the
	// count of shared secrets for an org and repo in the database.
	//
	// nolint: gosec // ignore false positive
	SelectSharedSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'shared'
AND org = ?
AND team = ?;
`

	// SelectOrgSecret represents a query to select a
	// secret for an org and name in the database.
	//
	// nolint: gosec // ignore false positive
	SelectOrgSecret = `
SELECT *
FROM secrets
WHERE type = 'org'
AND org = ?
AND name = ?
LIMIT 1;
`

	// SelectRepoSecret represents a query to select a
	// secret for an org, repo and name in the database.
	//
	// nolint: gosec // ignore false positive
	SelectRepoSecret = `
SELECT *
FROM secrets
WHERE type = 'repo'
AND org = ?
AND repo = ?
AND name = ?
LIMIT 1;
`

	// SelectSharedSecret represents a query to select a
	// secret for an org, team and name in the database.
	//
	// nolint: gosec // ignore false positive
	SelectSharedSecret = `
SELECT *
FROM secrets
WHERE type = 'shared'
AND org = ?
AND team = ?
AND name = ?
LIMIT 1;
`

	// DeleteSecret represents a query to
	// remove a secret from the database.
	//
	// nolint: gosec // ignore false positive
	DeleteSecret = `
DELETE
FROM secrets
WHERE id = ?;
`
)
