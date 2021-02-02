// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

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
AND org = $1
ORDER BY id DESC
LIMIT $2
OFFSET $3;
`

	// ListRepoSecrets represents a query to list all
	// secrets for a type, org and repo in the database.
	//
	// nolint: gosec // ignore false positive
	ListRepoSecrets = `
SELECT *
FROM secrets
WHERE type = 'repo'
AND org = $1
AND repo = $2
ORDER BY id DESC
LIMIT $3
OFFSET $4;
`

	// ListSharedSecrets represents a query to list all
	// secrets for a type, org and team in the database.
	//
	// nolint: gosec // ignore false positive
	ListSharedSecrets = `
SELECT *
FROM secrets
WHERE type = 'shared'
AND org = $1
AND team = $2
ORDER BY id DESC
LIMIT $3
OFFSET $4;
`

	// SelectOrgSecretsCount represents a query to select the
	// count of org secrets for an org in the database.
	//
	// nolint: gosec // ignore false positive
	SelectOrgSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'org'
AND org = $1;
`

	// SelectRepoSecretsCount represents a query to select the
	// count of repo secrets for an org and repo in the database.
	//
	// nolint: gosec // ignore false positive
	SelectRepoSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'repo'
AND org = $1
AND repo = $2;
`

	// SelectSharedSecretsCount represents a query to select the
	// count of shared secrets for an org and repo in the database.
	//
	// nolint: gosec // ignore false positive
	SelectSharedSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'shared'
AND org = $1
AND team = $2;
`

	// SelectOrgSecret represents a query to select a
	// secret for an org and name in the database.
	//
	// nolint: gosec // ignore false positive
	SelectOrgSecret = `
SELECT *
FROM secrets
WHERE type = 'org'
AND org = $1
AND name = $2
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
AND org = $1
AND repo = $2
AND name = $3
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
AND org = $1
AND team = $2
AND name = $3
LIMIT 1;
`

	// DeleteSecret represents a query to
	// remove a secret from the database.
	//
	// nolint: gosec // ignore false positive
	DeleteSecret = `
DELETE
FROM secrets
WHERE id = $1;
`
)

// createSecretService is a helper function to return
// a service for interacting with the secrets table.
func createSecretService() *Service {
	return &Service{
		List: map[string]string{
			"all":    ListSecrets,
			"org":    ListOrgSecrets,
			"repo":   ListRepoSecrets,
			"shared": ListSharedSecrets,
		},
		Select: map[string]string{
			"org":         SelectOrgSecret,
			"repo":        SelectRepoSecret,
			"shared":      SelectSharedSecret,
			"countOrg":    SelectOrgSecretsCount,
			"countRepo":   SelectRepoSecretsCount,
			"countShared": SelectSharedSecretsCount,
		},
		Delete: DeleteSecret,
	}
}
