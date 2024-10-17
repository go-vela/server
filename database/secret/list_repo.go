// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/types"
	"github.com/go-vela/types/constants"
)

// ListSecretsForRepo gets a list of secrets by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListSecretsForRepo(ctx context.Context, r *api.Repo, filters map[string]interface{}, page, perPage int) ([]*api.Secret, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"type": constants.SecretRepo,
	}).Tracef("listing secrets for repo %s", r.GetFullName())

	// variables to store query results and return values
	count := int64(0)
	s := new([]types.Secret)
	secrets := []*api.Secret{}

	// count the results
	count, err := e.CountSecretsForRepo(ctx, r, filters)
	if err != nil {
		return secrets, 0, err
	}

	// short-circuit if there are no results
	if count == 0 {
		return secrets, 0, nil
	}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err = e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretRepo).
		Where("org = ?", r.GetOrg()).
		Where("repo = ?", r.GetName()).
		Where(filters).
		Order("id DESC").
		Limit(perPage).
		Offset(offset).
		Find(&s).
		Error
	if err != nil {
		return nil, count, err
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

	return secrets, count, nil
}
