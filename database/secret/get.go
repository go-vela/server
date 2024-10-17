// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// GetSecret gets a secret by ID from the database.
func (e *engine) GetSecret(ctx context.Context, id int64) (*api.Secret, error) {
	e.logger.Tracef("getting secret %d", id)

	// variable to store query results
	s := new(types.Secret)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	err = s.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted secrets
		e.logger.Errorf("unable to decrypt secret %d: %v", id, err)
	}

	return s.ToAPI(), nil
}
