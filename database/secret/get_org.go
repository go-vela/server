// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetSecretForOrg gets a secret by org name from the database.
func (e *engine) GetSecretForOrg(ctx context.Context, org, name string) (*api.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":    org,
		"secret": name,
		"type":   constants.SecretOrg,
	}).Tracef("getting org secret %s/%s", org, name)

	// variable to store query results
	s := new(types.Secret)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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
	}

	return s.ToAPI(), nil
}
