// SPDX-License-Identifier: Apache-2.0

package favorite

import "context"

const (
	// CreateOrgNameIndex represents a query to create an
	// index on the repos table for the org and name columns.
	CreateRepoIndex = `
	CREATE INDEX idx_favorites_repo ON favorites (repo_id);
`

	CreateUserRepoIndex = `
	CREATE INDEX idx_favorites_user_position ON favorites (user_id, position);
`
)

// CreateFavoritesIndexes creates the indexes for the favorites table in the database.
func (e *Engine) CreateFavoritesIndexes(ctx context.Context) error {
	e.logger.Tracef("creating indexes for favorites table")

	// create the repo columns index for the favorites table
	if err := e.client.
		WithContext(ctx).
		Exec(CreateRepoIndex).Error; err != nil {
		return err
	}

	// create the user and position columns index for the favorites table
	if err := e.client.
		WithContext(ctx).
		Exec(CreateUserRepoIndex).Error; err != nil {
		return err
	}

	return nil
}
