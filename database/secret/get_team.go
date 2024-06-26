// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetSecretForTeam gets a secret by org and team name from the database.
func (e *engine) GetSecretForTeam(ctx context.Context, org, team, name string) (*library.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":    org,
		"team":   team,
		"secret": name,
		"type":   constants.SecretShared,
	}).Tracef("getting shared secret %s/%s/%s", org, team, name)

	// variable to store query results
	s := new(database.Secret)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretShared).
		Where("org = ?", org).
		Where("team = ?", team).
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
		e.logger.Errorf("unable to decrypt shared secret %s/%s/%s: %v", org, team, name, err)

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
