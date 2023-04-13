// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"strings"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// CountSecretsForTeam gets the count of secrets by org and team name from the database.
func (e *engine) CountSecretsForTeam(org, team string, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"team": team,
		"type": constants.SecretShared,
	}).Tracef("getting count of secrets for team %s/%s from the database", org, team)

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretShared).
		Where("org = ?", org).
		Where("team = ?", team).
		Where(filters).
		Count(&s).
		Error

	return s, err
}

// CountSecretsForTeams gets the count of secrets by teams within an org from the database.
func (e *engine) CountSecretsForTeams(org string, teams []string, filters map[string]interface{}) (int64, error) {
	// lower case team names for not case-sensitive values from the SCM i.e. GitHub
	//
	// iterate through the list of teams provided
	for index, team := range teams {
		// ensure the team name is lower case
		teams[index] = strings.ToLower(team)
	}

	e.logger.WithFields(logrus.Fields{
		"org":   org,
		"teams": teams,
		"type":  constants.SecretShared,
	}).Tracef("getting count of secrets for teams %s in org %s from the database", teams, org)

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretShared).
		Where("org = ?", org).
		Where("LOWER(team) IN (?)", teams).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
