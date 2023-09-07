// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListSecretsForOrg gets a list of secrets by org name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListSecretsForOrg(ctx context.Context, org string, filters map[string]interface{}, page, perPage int) ([]*library.Secret, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"type": constants.SecretOrg,
	}).Tracef("listing secrets for org %s from the database", org)

	// variables to store query results and return values
	count := int64(0)
	s := new([]database.Secret)
	secrets := []*library.Secret{}

	// count the results
	count, err := e.CountSecretsForOrg(ctx, org, filters)
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
		return nil, count, err
	}

	// iterate through all query results
	for _, secret := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		// decrypt the fields for the secret
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.Decrypt
		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			e.logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		// convert query result to library type
		//
		// https://pkg.go.dev/github.com/go-vela/types/database#Secret.ToLibrary
		secrets = append(secrets, tmp.ToLibrary())
	}

	return secrets, count, nil
}
