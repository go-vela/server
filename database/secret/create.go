// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateSecret creates a new secret in the database.
func (e *Engine) CreateSecret(ctx context.Context, s *api.Secret) (*api.Secret, error) {
	// handle the secret based off the type
	switch s.GetType() {
	case constants.SecretShared:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"team":   s.GetTeam(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("creating secret %s/%s/%s/%s", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName())
	default:
		e.logger.WithFields(logrus.Fields{
			"org":    s.GetOrg(),
			"repo":   s.GetRepo(),
			"secret": s.GetName(),
			"type":   s.GetType(),
		}).Tracef("creating secret %s/%s/%s/%s", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName())
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

	var result *api.Secret

	transactionErr := e.client.Transaction(func(tx *gorm.DB) error {
		// create secret record
		err = tx.
			WithContext(ctx).
			Table(constants.TableSecret).
			Create(secret.Nullify()).Error
		if err != nil {
			return err
		}

		// decrypt the fields for the secret to return
		err = secret.Decrypt(e.config.EncryptionKey)
		if err != nil {
			switch s.GetType() {
			case constants.SecretShared:
				return fmt.Errorf("unable to decrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetTeam(), s.GetName(), err)
			default:
				return fmt.Errorf("unable to decrypt secret %s/%s/%s/%s: %w", s.GetType(), s.GetOrg(), s.GetRepo(), s.GetName(), err)
			}
		}

		result = secret.ToAPI()
		result.SetRepoAllowlist(s.GetRepoAllowlist())

		err = InsertAllowlist(ctx, tx, result)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionErr != nil {
		return nil, transactionErr
	}

	return result, nil
}
