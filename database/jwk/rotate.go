// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"database/sql"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// RotateKeys removes all inactive keys and sets active keys to inactive.
func (e *Engine) RotateKeys(ctx context.Context) error {
	e.logger.Trace("rotating jwks")

	k := types.JWK{}

	// remove inactive keys
	err := e.client.
		WithContext(ctx).
		Table(constants.TableJWK).
		Where("active = ?", false).
		Delete(&k).
		Error
	if err != nil {
		return err
	}

	// set active keys to inactive
	err = e.client.
		WithContext(ctx).
		Table(constants.TableJWK).
		Where("active = ?", true).
		Update("active", sql.NullBool{Bool: false, Valid: true}).
		Error
	if err != nil {
		return err
	}

	return nil
}
