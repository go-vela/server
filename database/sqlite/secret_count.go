// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/database/sqlite/dml"
	"github.com/go-vela/types/constants"
)

// GetTypeSecretCount gets a count of secrets by type,
// owner, and name (repo or team) from the database.
func (c *client) GetTypeSecretCount(t, o, n string, teams []string) (int64, error) {
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
	c.Logger.WithFields(fields).Tracef("getting count of %s secrets for %s/%s from the database", t, o, n)

	var err error

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Sqlite.
			Table(constants.TableSecret).
			Raw(dml.SelectOrgSecretsCount, o).
			Pluck("count", &s).Error
	case constants.SecretRepo:
		err = c.Sqlite.
			Table(constants.TableSecret).
			Raw(dml.SelectRepoSecretsCount, o, n).
			Pluck("count", &s).Error
	case constants.SecretShared:
		if n == "*" {
			// GitHub teams are not case-sensitive, the DB is lowercase everything for matching
			var lowerTeams []string
			for _, t := range teams {
				lowerTeams = append(lowerTeams, strings.ToLower(t))
			}
			err = c.Sqlite.
				Table(constants.TableSecret).
				Select("count(*)").
				Where("type = 'shared' AND org = ?", o).
				Where("LOWER(team) IN (?)", lowerTeams).
				Pluck("count", &s).Error
		} else {
			err = c.Sqlite.
				Table(constants.TableSecret).
				Raw(dml.SelectSharedSecretsCount, o, n).
				Pluck("count", &s).Error
		}
	}

	return s, err
}
