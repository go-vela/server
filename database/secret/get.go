// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetSecret gets a secret by ID from the database.
func (e *engine) GetSecret(ctx context.Context, id int64) (*library.Secret, error) {
	e.logger.Tracef("getting secret %d from the database", id)

	// variable to store query results
	s := new(database.Secret)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("id = ?", id).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// decrypt the fields for the secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
	err = s.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted secrets
		e.logger.Errorf("unable to decrypt secret %d: %v", id, err)

		// return the unencrypted secret
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.ToLibrary
		return s.ToLibrary(), nil
	}

	// return the decrypted secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.ToLibrary
	return s.ToLibrary(), nil
}
