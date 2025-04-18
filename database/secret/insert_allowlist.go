// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"gorm.io/gorm/clause"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// InsertAllowlist adds allowlist entries in the database.
func (e *Engine) InsertAllowlist(ctx context.Context, s *api.Secret) error {
	// only insert when there is an allowlist
	if len(s.GetRepoAllowlist()) == 0 {
		return nil
	}

	entries := []*types.SecretAllowlist{}

	for _, r := range s.GetRepoAllowlist() {
		entries = append(entries, types.SecretAllowlistFromAPI(s, r))
	}

	// upsert allowlist
	return e.client.
		WithContext(ctx).
		Table(constants.TableSecretRepoAllowlist).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(entries).Error
}
