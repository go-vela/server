// SPDX-License-Identifier: Apache-2.0

package secret

import "context"

const (
	// CreateTypeOrgRepo represents a query to create an
	// index on the secrets table for the type, org and repo columns.
	CreateTypeOrgRepo = `
CREATE INDEX
IF NOT EXISTS
secrets_type_org_repo
ON secrets (type, org, repo);
`
	// CreateTypeOrgTeam represents a query to create an
	// index on the secrets table for the type, org and team columns.
	CreateTypeOrgTeam = `
CREATE INDEX
IF NOT EXISTS
secrets_type_org_team
ON secrets (type, org, team);
`
	// CreateTypeOrg represents a query to create an
	// index on the secrets table for the type, and org columns.
	CreateTypeOrg = `
CREATE INDEX
IF NOT EXISTS
secrets_type_org
ON secrets (type, org);
`

	// CreateSecretID represents a query to create an
	// index on the secret_repo_allowlist tabe for the secret_id column.
	//nolint:gosec // not credentials
	CreateSecretID = `
CREATE INDEX
IF NOT EXISTS
secret_repo_allowlist_secret_id
ON secret_repo_allowlist (secret_id)
`
)

// CreateSecretIndexes creates the indexes for the secrets table in the database.
func (e *Engine) CreateSecretIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for secrets table")

	// create the secret allowlist secret_id index for the secret_repo_allowlist table
	err := e.client.
		WithContext(ctx).
		Exec(CreateSecretID).Error
	if err != nil {
		return err
	}

	// create the type, org and repo columns index for the secrets table
	err = e.client.
		WithContext(ctx).
		Exec(CreateTypeOrgRepo).Error
	if err != nil {
		return err
	}

	// create the type, org and team columns index for the secrets table
	err = e.client.
		WithContext(ctx).
		Exec(CreateTypeOrgTeam).Error
	if err != nil {
		return err
	}

	// create the type and org columns index for the secrets table
	return e.client.
		WithContext(ctx).
		Exec(CreateTypeOrg).Error
}
