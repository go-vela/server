// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// GetTypeSecretCount gets a count of secrets by type,
// owner, and name (repo or team) from the database.
func (c *client) GetTypeSecretCount(t, o, n string) (int64, error) {
	logrus.Tracef("getting count of %s secrets for %s/%s from the database", t, o, n)

	var err error

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	switch t {
	case constants.SecretOrg:
		err = c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.SelectOrgSecretsCount, o).
			Pluck("count", &s).Error
	case constants.SecretRepo:
		err = c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.SelectRepoSecretsCount, o, n).
			Pluck("count", &s).Error
	case constants.SecretShared:
		err = c.Postgres.
			Table(constants.TableSecret).
			Raw(dml.SelectSharedSecretsCount, o, n).
			Pluck("count", &s).Error
	}

	return s, err
}
