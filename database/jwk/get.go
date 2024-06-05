// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"

	"github.com/lestrrat-go/jwx/jwk"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetActiveJWK gets a JWK by UUID (kid) from the database if active.
func (e *engine) GetActiveJWK(_ context.Context, id string) (jwk.RSAPublicKey, error) {
	e.logger.Tracef("getting key %s from the database", id)

	// variable to store query results
	j := new(types.JWK)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableJWK).
		Where("id = ?", id).
		Where("active = ?", true).
		Take(j).
		Error
	if err != nil {
		return j.ToAPI(), err
	}

	return j.ToAPI(), nil
}
