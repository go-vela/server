// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// CountSecretsForRepo gets the count of secrets by org and repo name from the database.
func (e *engine) CountSecretsForRepo(ctx context.Context, r *library.Repo, filters map[string]interface{}) (int64, error) {
	e.logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"type": constants.SecretRepo,
	}).Tracef("getting count of secrets for repo %s from the database", r.GetFullName())

	// variable to store query results
	var s int64

	// send query to the database and store result in variable
	err := e.client.
		Table(constants.TableSecret).
		Where("type = ?", constants.SecretRepo).
		Where("org = ?", r.GetOrg()).
		Where("repo = ?", r.GetName()).
		Where(filters).
		Count(&s).
		Error

	return s, err
}
