// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with create.go
package secret

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateSecret updates an existing secret in the database.
func (e *engine) UpdateSecret(ctx context.Context, s *api.Secret) (*api.Secret, error) {
	// handle the secret based off the type
	switch s.GetType() {
	case constants.SecretShared:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"team":   s.GetTeam(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("updating secret %s/%s/%s/%s", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName())
	default:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"repo":   s.GetRepo(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("updating secret %s/%s/%s/%s", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName())
	}

	secret := types.SecretFromAPI(s)

	err := secret.Validate()
	if err != nil {
		return nil, err
	}

	err = secret.Encrypt(e.config.EncryptionKey)
	if err != nil {
		switch s.GetType() {
		case constants.SecretShared:
			return nil, fmt.Errorf("unable to encrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName(), err)
		default:
			return nil, fmt.Errorf("unable to encrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName(), err)
		}
	}

	err = e.client.
		WithContext(ctx).
		Table(constants.TableSecret).
		Save(secret.Nullify()).Error
	if err != nil {
		return nil, err
	}

	err = secret.Decrypt(e.config.EncryptionKey)
	if err != nil {
		switch s.GetType() {
		case constants.SecretShared:
			return nil, fmt.Errorf("unable to decrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName(), err)
		default:
			return nil, fmt.Errorf("unable to decrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName(), err)
		}
	}

	return secret.ToAPI(), nil
}
