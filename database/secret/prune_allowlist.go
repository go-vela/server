// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
	"gorm.io/gorm"
)

// PruneAllowlist deletes any allowlist record from the database that belongs to the secret but is not in the active allowlist.
func PruneAllowlist(ctx context.Context, tx *gorm.DB, s *api.Secret) error {
	// send query to the database
	return tx.
		WithContext(ctx).
		Table(constants.TableSecretRepoAllowlist).
		Where("secret_id = ?", s.GetID()).
		Where("repo NOT IN (?)", s.GetRepoAllowlist()).
		Delete(&types.SecretAllowlist{}).
		Error
}
