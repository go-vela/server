// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetSecretForRepo gets a secret by org and repo name from the database.
func (e *engine) GetSecretForRepo(ctx context.Context, name string, r *api.Repo) (*library.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":    r.GetOrg(),
		"repo":   r.GetName(),
		"secret": name,
		"type":   constants.SecretRepo,
	}).Tracef("getting repo secret %s/%s", r.GetFullName(), name)

	// variable to store query results
	s := new(database.Secret)

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretRepo).
		Where("org = ?", r.GetOrg()).
		Where("repo = ?", r.GetName()).
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
		e.logger.Errorf("unable to decrypt repo secret %s/%s: %v", r.GetFullName(), name, err)

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
