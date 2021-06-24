// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"errors"
	"fmt"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

// GetSecret gets a secret by type, org, name (repo or team) and secret name from the database.
func (c *client) GetSecret(t, o, n, secretName string) (*library.Secret, error) {
	logrus.Tracef("getting %s secret %s for %s/%s from the database", t, secretName, o, n)

	var err error

	// variable to store query results
	s := new(database.Secret)

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		result := c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.SelectOrgSecret, o, secretName).
			Scan(s)

		// check if the query returned a record not found error or no rows were returned
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}

		err = result.Error
	case constants.SecretRepo:
		result := c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.SelectRepoSecret, o, n, secretName).
			Scan(s)

		// check if the query returned a record not found error or no rows were returned
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}

		err = result.Error
	case constants.SecretShared:
		result := c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.SelectSharedSecret, o, n, secretName).
			Scan(s)

		// check if the query returned a record not found error or no rows were returned
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result.RowsAffected == 0 {
			return nil, gorm.ErrRecordNotFound
		}

		err = result.Error
	}
	if err != nil {
		return nil, err
	}

	// decrypt the value for the secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
	err = s.Decrypt(c.config.EncryptionKey)
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted secrets
		logrus.Errorf("unable to decrypt %s secret %s for %s/%s: %v", t, secretName, o, n, err)

		// return the unencrypted secret
		return s.ToLibrary(), nil
	}

	// return the decrypted secret
	return s.ToLibrary(), nil
}

// CreateSecret creates a new secret in the database.
func (c *client) CreateSecret(s *library.Secret) error {
	logrus.Tracef("creating %s secret %s in the database", s.GetType(), s.GetName())

	// cast to database type
	secret := database.SecretFromLibrary(s)

	// validate the necessary fields are populated
	err := secret.Validate()
	if err != nil {
		return err
	}

	// encrypt the value for the secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Encrypt
	err = secret.Encrypt(c.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt secret %s: %v", s.GetName(), err)
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableSecret).
		Create(secret.Nullify()).Error
}

// UpdateSecret updates a secret in the database.
func (c *client) UpdateSecret(s *library.Secret) error {
	logrus.Tracef("updating %s secret %s in the database", s.GetType(), s.GetName())

	// cast to database type
	secret := database.SecretFromLibrary(s)

	// validate the necessary fields are populated
	err := secret.Validate()
	if err != nil {
		return err
	}

	// encrypt the value for the secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Encrypt
	err = secret.Encrypt(c.config.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt secret %s: %v", s.GetName(), err)
	}

	// send query to the database
	return c.Postgres.
		Table(constants.TableSecret).
		Save(secret.Nullify()).Error
}

// DeleteSecret deletes a secret by unique ID from the database.
func (c *client) DeleteSecret(id int64) error {
	logrus.Tracef("Deleting secret %d from the database", id)

	// send query to the database
	return c.Postgres.
		Table(constants.TableSecret).
		Exec(dml.DeleteSecret, id).Error
}
