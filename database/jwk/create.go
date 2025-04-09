// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateJWK creates a new JWK in the database.
func (e *Engine) CreateJWK(ctx context.Context, j jwk.RSAPublicKey) error {
	logKeyID, ok := j.KeyID()
	if !ok {
		return fmt.Errorf("unable to create JWK: no key provided")
	}

	e.logger.WithFields(logrus.Fields{
		"jwk": logKeyID,
	}).Tracef("creating key %s", logKeyID)

	key := types.JWKFromAPI(j)
	key.Active = sql.NullBool{Bool: true, Valid: true}

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableJWK).
		Create(key).Error
}
