// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteSecret deletes an existing secret from the database.
func (e *Engine) DeleteSecret(ctx context.Context, s *api.Secret) error {
	// handle the secret based off the type
	//
	//nolint:dupl // ignore similar code with update.go
	switch s.GetType() {
	case constants.SecretShared:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"team":   s.GetTeam(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("deleting secret %s/%s/%s/%s", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName())
	default:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"repo":   s.GetRepo(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("deleting secret %s/%s/%s/%s", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName())
	}

	secret := types.SecretFromAPI(s)

	// send query to the database
	err := e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Delete(secret).
		Error
	if err != nil {
		return err
	}

	// empty allowlist
	s.SetRepoAllowlist([]string{})

	return e.PruneAllowlist(ctx, s)
}
