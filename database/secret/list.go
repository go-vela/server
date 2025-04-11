// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListSecrets gets a list of all secrets from the database.
func (e *Engine) ListSecrets(ctx context.Context) ([]*api.Secret, error) {
	e.logger.Trace("listing all secrets")

	// variables to store query results and return value
	s := new([]types.Secret)
	secrets := []*api.Secret{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Find(&s).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, secret := range *s {
		err = secret.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			e.logger.Errorf("unable to decrypt secret %d: %v", secret.ID.Int64, err)
		}

		secrets = append(secrets, secret.ToAPI())
	}

	return secrets, nil
}
