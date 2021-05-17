// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// GetSecretList gets a list of all secrets from the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) GetSecretList() ([]*library.Secret, error) {
	logrus.Tracef("listing secrets from the database")

	// variable to store query results
	s := new([]database.Secret)

	// send query to the database and store result in variable
	err := c.Postgres.
		Table(constants.TableSecret).
		Raw(dml.ListSecrets).
		Scan(s).Error
	if err != nil {
		return nil, err
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
		err = tmp.Decrypt(c.config.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			logrus.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, nil
}

// GetTypeSecretList gets a list of secrets by type,
// owner, and name (repo or team) from the database.
func (c *client) GetTypeSecretList(t, o, n string, page, perPage int) ([]*library.Secret, error) {
	logrus.Tracef("listing %s secrets for %s/%s from the database", t, o, n)

	var err error

	// variable to store query results
	s := new([]database.Secret)
	// calculate offset for pagination through results
	offset := (perPage * (page - 1))

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.ListOrgSecrets, o, perPage, offset).
			Scan(s).Error
	case constants.SecretRepo:
		err = c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.ListRepoSecrets, o, n, perPage, offset).
			Scan(s).Error
	case constants.SecretShared:
		err = c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.ListSharedSecrets, o, n, perPage, offset).
			Scan(s).Error
	}
	if err != nil {
		return nil, err
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
		err = tmp.Decrypt(c.config.EncryptionKey)
		if err != nil {
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			logrus.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, nil
}
