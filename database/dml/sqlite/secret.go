// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

const (
	// ListSecrets represents a query to
	// list all secrets in the database.
	ListSecrets = `
SELECT *
FROM secrets;
`

	// ListOrgSecrets represents a query to list all
	// secrets for a type and org in the database.
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
	SelectOrgSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'org'
AND org = ?;
`

	// SelectRepoSecretsCount represents a query to select the
	// count of repo secrets for an org and repo in the database.
	SelectRepoSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'repo'
AND org = ?
AND repo = ?;
`

	// SelectSharedSecretsCount represents a query to select the
	// count of shared secrets for an org and repo in the database.
	SelectSharedSecretsCount = `
SELECT count(*) as count
FROM secrets
WHERE type = 'shared'
AND org = ?
AND team = ?;
`

	// SelectOrgSecret represents a query to select a
	// secret for an org and name in the database.
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
	DeleteSecret = `
DELETE
FROM secrets
WHERE id = ?;
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
