// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetSecret gets a secret by type, org, name (repo or team) and secret name from the database.
func (c *client) GetSecret(t, o, n, secretName string) (*library.Secret, error) {
	logrus.Tracef("Getting %s secret %s for %s/%s from the database", t, secretName, o, n)

	var err error

	// variable to store query results
	s := new(database.Secret)

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.Select["org"], o, secretName).
			Scan(s).Error
	case constants.SecretRepo:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.Select["repo"], o, n, secretName).
			Scan(s).Error
	case constants.SecretShared:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.Select["shared"], o, n, secretName).
			Scan(s).Error
	}

	// decrypt the value for the secret
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
	err = s.Decrypt(c.EncryptionKey)
	if err != nil {
		// ensures that the change is backwards compatible
		// by logging the error instead of returning it
		// which allows us to fetch unencrypted secrets
		logrus.Errorf("unable to decrypt %s secret %s for %s/%s: %v", t, secretName, o, n, err)

		// return the unencrypted secret
		return s.ToLibrary(), nil
	}

	// return the decrypted log
	return s.ToLibrary(), err
}

// GetSecretList gets a list of all secrets from the database.
func (c *client) GetSecretList() ([]*library.Secret, error) {
	logrus.Tracef("Listing secrets from the database")

	// variable to store query results
	s := new([]database.Secret)

	// send query to the database and store result in variable
	err := c.Database.
		Table(constants.TableSecret).
		Raw(c.DML.SecretService.List["all"]).
		Scan(s).Error

	// variable we want to return
	secrets := []*library.Secret{}
	// iterate through all query results
	for _, secret := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		// decrypt the value for the secret
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
		err = tmp.Decrypt(c.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			logrus.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, err
}

// GetTypeSecretList gets a list of secrets by type,
// owner, and name (repo or team) from the database.
func (c *client) GetTypeSecretList(t, o, n string, page, perPage int) ([]*library.Secret, error) {
	logrus.Tracef("Listing %s secrets for %s/%s from the database", t, o, n)

	var err error

	// variable to store query results
	s := new([]database.Secret)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.List["org"], o, perPage, offset).
			Scan(s).Error
	case constants.SecretRepo:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.List["repo"], o, n, perPage, offset).
			Scan(s).Error
	case constants.SecretShared:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.List["shared"], o, n, perPage, offset).
			Scan(s).Error
	}

	// variable we want to return
	secrets := []*library.Secret{}
	// iterate through all query results
	for _, secret := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		// decrypt the value for the secret
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
		err = tmp.Decrypt(c.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			logrus.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, err
}

// GetTypeSecretCount gets a count of secrets by type,
// owner, and name (repo or team) from the database.
func (c *client) GetTypeSecretCount(t, o, n string) (int64, error) {
	logrus.Tracef("Counting %s secrets for %s/%s from the database", t, o, n)

	var err error

	// variable to store query results
	var r []int64

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.Select["countOrg"], o).
			Pluck("count", &r).Error
	case constants.SecretRepo:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.Select["countRepo"], o, n).
			Pluck("count", &r).Error
	case constants.SecretShared:
		err = c.Database.
			Table(constants.TableSecret).
			Raw(c.DML.SecretService.Select["countShared"], o, n).
			Pluck("count", &r).Error
	}

	// return 0 if no result was found
	if len(r) == 0 {
		return 0, err
	}

	return r[0], err
}

// CreateSecret creates a new secret in the database.
func (c *client) CreateSecret(s *library.Secret) error {
	logrus.Tracef("Creating secret %s in the database", s.GetName())

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
	err = secret.Encrypt(c.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt secret %s: %v", s.GetName(), err)
	}

	// send query to the database
	return c.Database.
		Table(constants.TableSecret).
		Create(secret).Error
}

// UpdateSecret updates a secret in the database.
func (c *client) UpdateSecret(s *library.Secret) error {
	logrus.Tracef("Updating secret %s in the database", s.GetName())

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
	err = secret.Encrypt(c.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to encrypt secret %s: %v", s.GetName(), err)
	}

	// send query to the database
	return c.Database.
		Table(constants.TableSecret).
		Where("id = ?", s.GetID()).
		Update(secret).Error
}

// DeleteSecret deletes a secret by unique ID from the database.
func (c *client) DeleteSecret(id int64) error {
	logrus.Tracef("Deleting secret %d from the database", id)

	// send query to the database
	return c.Database.
		Table(constants.TableSecret).
		Exec(c.DML.SecretService.Delete, id).Error
}
