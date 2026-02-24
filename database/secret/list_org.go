// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListSecretsForOrg gets a list of secrets by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *Engine) ListSecretsForOrg(ctx context.Context, org string, filters map[string]any, page, perPage int) ([]*api.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"type": constants.SecretOrg,
	}).Tracef("listing secrets for org %s", org)

	// variables to store query results and return values
	s := new([]types.Secret)
	secrets := []*api.Secret{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretOrg).
		Where("org = ?", org).
		Where(filters).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
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

	return e.FillSecretsAllowlists(ctx, secrets)
}
