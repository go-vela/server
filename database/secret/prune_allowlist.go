// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// PruneAllowlist deletes any allowlist record from the database that belongs to the secret but is not in the active allowlist.
func PruneAllowlist(ctx context.Context, tx *gorm.DB, s *api.Secret) error {
	// if allowlist is 0, do not use NOT IN clause
	if len(s.GetRepoAllowlist()) == 0 {
		return tx.
			WithContext(ctx).
			Table(constants.TableSecretRepoAllowlist).
			Where("secret_id = ?", s.GetID()).
			Delete(&types.SecretAllowlist{}).
			Error
	}

	return tx.
		WithContext(ctx).
		Table(constants.TableSecretRepoAllowlist).
		Where("secret_id = ?", s.GetID()).
		Where("repo NOT IN (?)", s.GetRepoAllowlist()).
		Delete(&types.SecretAllowlist{}).
		Error
}
