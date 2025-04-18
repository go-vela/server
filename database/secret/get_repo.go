// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetSecretForRepo gets a secret by org and repo name from the database.
func (e *Engine) GetSecretForRepo(ctx context.Context, name string, r *api.Repo) (*api.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":    r.GetOrg(),
		"repo":   r.GetName(),
		"secret": name,
		"type":   constants.SecretRepo,
	}).Tracef("getting repo secret %s/%s", r.GetFullName(), name)

	// variable to store query results
	s := new(types.Secret)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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

	err = s.Decrypt(e.config.EncryptionKey)
	if err != nil {
		// TODO: remove backwards compatibility before 1.x.x release
		//
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted secrets
		e.logger.Errorf("unable to decrypt repo secret %s/%s: %v", r.GetFullName(), name, err)
	}

	return e.FillSecretAllowlist(ctx, s.ToAPI())
}
