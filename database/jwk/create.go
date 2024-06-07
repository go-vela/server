// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"database/sql"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateJWK creates a new JWK in the database.
func (e *engine) CreateJWK(_ context.Context, j jwk.RSAPublicKey) error {
	e.logger.WithFields(logrus.Fields{
		"jwk": j.KeyID(),
	}).Tracef("creating key %s in the database", j.KeyID())

	key := types.JWKFromAPI(j)
	key.Active = sql.NullBool{Bool: true, Valid: true}

	// send query to the database
	return e.client.Table(constants.TableJWK).Create(key).Error
}
