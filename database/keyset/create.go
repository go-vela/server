// SPDX-License-Identifier: Apache-2.0

package keyset

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateKeySet creates a new dashboard in the database.
func (e *engine) CreateKeySet(ctx context.Context, j api.JWK) error {
	e.logger.WithFields(logrus.Fields{
		"jwk": j.Kid,
	}).Tracef("creating key %s in the database", j.Kid)

	key := types.JWKFromAPI(j)
	key.Active = sql.NullBool{Bool: true, Valid: true}

	err := key.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.Table(constants.TableKeySet).Create(key).Error

}
