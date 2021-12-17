// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
)

// GetSecretList gets a list of all secrets from the database.
//
// nolint: dupl // ignore false positive of duplicate code
func (c *client) GetSecretList() ([]*library.Secret, error) {
	c.Logger.Tracef("listing secrets from the database")

	// variable to store query results
	s := new([]database.Secret)

	// send query to the database and store result in variable
	err := c.Sqlite.
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
			c.Logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, nil
}

// GetTypeSecretList gets a list of secrets by type,
// owner, and name (repo or team) from the database.
//
// nolint: lll // ignore long line length
func (c *client) GetTypeSecretList(t, o, n string, page, perPage int, teams []string) ([]*library.Secret, error) {
	// create log fields from secret metadata
	fields := logrus.Fields{
		"org":  o,
		"repo": n,
		"type": t,
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields from secret metadata
		fields = logrus.Fields{
			"org":  o,
			"team": n,
			"type": t,
		}
	}

	// nolint: lll // ignore long line length due to parameters
	c.Logger.WithFields(fields).Tracef("listing %s secrets for %s/%s from the database", t, o, n)

	var err error

	// variable to store query results
	s := new([]database.Secret)
	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Sqlite.
			Table(constants.TableSecret).
			Raw(dml.ListOrgSecrets, o, perPage, offset).
			Scan(s).Error
	case constants.SecretRepo:
		err = c.Sqlite.
			Table(constants.TableSecret).
			Raw(dml.ListRepoSecrets, o, n, perPage, offset).
			Scan(s).Error
	case constants.SecretShared:
		if n == "*" {
			// GitHub teams are not case-sensitive, the DB is lowercase everything for matching
			var lowerTeams []string
			for _, t := range teams {
				lowerTeams = append(lowerTeams, strings.ToLower(t))
			}
			err = c.Sqlite.
				Table(constants.TableSecret).
				Where("type = 'shared' AND org = ?", o).
				Where("LOWER(team) IN (?)", lowerTeams).
				Order("id DESC").
				Limit(perPage).
				Offset(offset).
				Scan(s).Error
		} else {
			err = c.Sqlite.
				Table(constants.TableSecret).
				Raw(dml.ListSharedSecrets, o, n, perPage, offset).
				Scan(s).Error
		}
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
			c.Logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, nil
}
