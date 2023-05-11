// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// ListSecretsForTeam gets a list of secrets by org and team name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *engine) ListSecretsForTeam(org, team string, filters map[string]interface{}, page, perPage int) ([]*library.Secret, int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"team": team,
		"type": constants.SecretShared,
	}).Tracef("listing secrets for team %s/%s from the database", org, team)

	// variables to store query results and return values
	count := int64(0)
	s := new([]database.Secret)
	secrets := []*library.Secret{}

	// count the results
	count, err := e.CountSecretsForTeam(org, team, filters)
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
		Where("type = ?", constants.SecretShared).
		Where("org = ?", org).
		Where("team = ?", team).
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

// ListSecretsForTeams gets a list of secrets by teams within an org from the database.
func (e *engine) ListSecretsForTeams(org string, teams []string, filters map[string]interface{}, page, perPage int) ([]*library.Secret, int64, error) {
	// iterate through the list of teams provided
	for index, team := range teams {
		// ensure the team name is lower case
		teams[index] = strings.ToLower(team)
	}

	e.logger.WithFields(logrus.Fields{
		"org":   org,
		"teams": teams,
		"type":  constants.SecretShared,
	}).Tracef("listing secrets for teams %s in org %s from the database", teams, org)

	// variables to store query results and return values
	count := int64(0)
	s := new([]database.Secret)
	secrets := []*library.Secret{}

	// count the results
	count, err := e.CountSecretsForTeams(org, teams, filters)
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
		Where("type = ?", constants.SecretShared).
		Where("org = ?", org).
		Where("LOWER(team) IN (?)", teams).
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
