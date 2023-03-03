// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

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
func (e *engine) CreateSecretIndexes() error {
	e.logger.Tracef("creating indexes for secrets table in the database")

	// create the type, org and repo columns index for the secrets table
	err := e.client.Exec(CreateTypeOrgRepo).Error
	if err != nil {
		return err
	}

	// create the type, org and team columns index for the secrets table
	err = e.client.Exec(CreateTypeOrgTeam).Error
	if err != nil {
		return err
	}

	// create the type and org columns index for the secrets table
	return e.client.Exec(CreateTypeOrg).Error
}