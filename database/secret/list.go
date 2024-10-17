// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListSecrets gets a list of all secrets from the database.
func (e *engine) ListSecrets(ctx context.Context) ([]*api.Secret, error) {
	e.logger.Trace("listing all secrets")

	// variables to store query results and return value
	count := int64(0)
	s := new([]types.Secret)
	secrets := []*api.Secret{}

	// count the results
	count, err := e.CountSecrets(ctx)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return secrets, nil
	}

	// send query to the database and store result in variable
	err = e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Find(&s).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, secret := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			e.logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		secrets = append(secrets, tmp.ToAPI())
	}

	return secrets, nil
}
