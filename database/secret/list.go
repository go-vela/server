// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// ListSecrets gets a list of all secrets from the database.
func (e *engine) ListSecrets() ([]*library.Secret, error) {
	e.logger.Trace("listing all secrets from the database")

	// variables to store query results and return value
	count := int64(0)
	s := new([]database.Secret)
	secrets := []*library.Secret{}

	// count the results
	count, err := e.CountSecrets()
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return secrets, nil
	}

	// send query to the database and store result in variable
	err = e.client.
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

		// decrypt the fields for the secret
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			e.logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.ToLibrary
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, nil
}
