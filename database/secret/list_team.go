// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// ListSecretsForTeam gets a list of secrets by org and team name from the database.
//
//nolint:lll // ignore long line length due to variable names
func (e *Engine) ListSecretsForTeam(ctx context.Context, org, team string, filters map[string]interface{}, page, perPage int) ([]*api.Secret, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"team": team,
		"type": constants.SecretShared,
	}).Tracef("listing secrets for team %s/%s", org, team)

	// variables to store query results and return values
	s := new([]types.Secret)
	secrets := []*api.Secret{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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
		return nil, err
	}

	// iterate through all query results
	for _, secret := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			e.logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		secrets = append(secrets, tmp.ToAPI())
	}

	return secrets, nil
}

// ListSecretsForTeams gets a list of secrets by teams within an org from the database.
func (e *Engine) ListSecretsForTeams(ctx context.Context, org string, teams []string, filters map[string]interface{}, page, perPage int) ([]*api.Secret, error) {
	// iterate through the list of teams provided
	for index, team := range teams {
		// ensure the team name is lower case
		teams[index] = strings.ToLower(team)
	}

	e.logger.WithFields(logrus.Fields{
		"org":   org,
		"teams": teams,
		"type":  constants.SecretShared,
	}).Tracef("listing secrets for teams %s in org %s", teams, org)

	// variables to store query results and return values
	s := new([]types.Secret)
	secrets := []*api.Secret{}

	// calculate offset for pagination through results
	offset := perPage * (page - 1)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
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
		return nil, err
	}

	// iterate through all query results
	for _, secret := range *s {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := secret

		err = tmp.Decrypt(e.config.EncryptionKey)
		if err != nil {
			// TODO: remove backwards compatibility before 1.x.x release
			//
			// ensures that the change is backwards compatible
			// by logging the error instead of returning it
			// which allows us to fetch unencrypted secrets
			e.logger.Errorf("unable to decrypt secret %d: %v", tmp.ID.Int64, err)
		}

		secrets = append(secrets, tmp.ToAPI())
	}

	return secrets, nil
}
