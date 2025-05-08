// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// FillSecretAllowlist gets a secret allowlist by secret id.
func (e *Engine) FillSecretAllowlist(ctx context.Context, s *api.Secret) (*api.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"secret_id": s.GetID(),
	}).Tracef("getting allowlist for secret %d", s.GetID())

	// variable to store query results
	allowlist := new([]types.SecretAllowlist)

	result := []string{}

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSecretRepoAllowlist).
		Where("secret_id = ?", s.GetID()).
		Find(allowlist).
		Error
	if err != nil {
		return nil, err
	}

	// iterate through all query results
	for _, record := range *allowlist {
		tmp := record

		result = append(result, tmp.Repo.String)
	}

	s.SetRepoAllowlist(result)

	return s, nil
}
