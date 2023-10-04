// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// GetSecretForOrg gets a secret by org name from the database.
func (e *engine) GetSecretForOrg(ctx context.Context, org, name string) (*library.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":    org,
		"secret": name,
		"type":   constants.SecretOrg,
	}).Tracef("getting org secret %s/%s from the database", org, name)

	// variable to store query results
	s := new(database.Secret)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretOrg).
		Where("org = ?", org).
		Where("name = ?", name).
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
		e.logger.Errorf("unable to decrypt org secret %s/%s: %v", org, name, err)

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
