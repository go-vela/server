// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"database/sql"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// RotateKeys removes all inactive keys and sets active keys to inactive.
func (e *engine) RotateKeys(_ context.Context) error {
	e.logger.Trace("rotating keysets in the database")

	k := types.JWK{}

	// remove inactive keys
	err := e.client.
		Table(constants.TableKeySet).
		Where("active = ?", false).
		Delete(&k).
		Error
	if err != nil {
		return err
	}

	// set active keys to inactive
	err = e.client.
		Table(constants.TableKeySet).
		Where("active = ?", true).
		Update("active", sql.NullBool{Bool: false, Valid: true}).
		Error
	if err != nil {
		return err
	}

	return nil
}
