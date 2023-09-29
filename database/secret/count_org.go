// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// CountSecretsForOrg gets the count of secrets by org name from the database.
func (e *engine) CountSecretsForOrg(ctx context.Context, org string, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  org,
		"type": constants.SecretOrg,
	}).Tracef("getting count of secrets for org %s from the database", org)

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretOrg).
		Where("org = ?", org).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
