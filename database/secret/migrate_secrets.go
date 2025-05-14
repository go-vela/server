// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

// MigrateSecrets updates repository secrets and allowlist records for a repo name change.
func (e *Engine) MigrateSecrets(ctx context.Context, oldOrg, oldRepo, newOrg, newRepo string) error {
	return e.client.Transaction(func(tx *gorm.DB) error {
		// set org and repo to new values for repo type secrets
		err := tx.
			WithContext(ctx).
			Table(constants.TableSecret).
			Where("type = ?", constants.SecretRepo).
			Where("org = ?", oldOrg).
			Where("repo = ?", oldRepo).
			Updates(map[string]any{
				"org":  newOrg,
				"repo": newRepo,
			}).Error
		if err != nil {
			return err
		}

		// set allowlist records involving old repo to use the new repo for all secret types with allowlists
		err = tx.
			WithContext(ctx).
			Table(constants.TableSecretRepoAllowlist).
			Where("repo = ?", fmt.Sprintf("%s/%s", oldOrg, oldRepo)).
			Updates(map[string]any{
				"repo": fmt.Sprintf("%s/%s", newOrg, newRepo),
			}).Error
		if err != nil {
			return err
		}

		return nil
	})
}
