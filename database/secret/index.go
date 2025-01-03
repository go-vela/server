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
)

// CreateSecretIndexes creates the indexes for the secrets table in the database.
func (e *engine) CreateSecretIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for secrets table")

	// create the type, org and repo columns index for the secrets table
	err := e.client.
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
