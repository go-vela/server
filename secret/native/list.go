// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// List captures a list of secrets.
func (c *client) List(sType, org, name string, page, perPage int, teams []string) ([]*library.Secret, int64, error) {
	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"type": sType,
		}).Tracef("listing native %s secrets for %s", sType, org)

		// capture the count of org secrets from the native service
		return c.Database.ListSecretsForOrg(org, nil, page, perPage)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"repo": name,
			"type": sType,
		}).Tracef("listing native %s secrets for %s/%s", sType, org, name)

		// create the repo with the information available
		r := new(library.Repo)
		r.SetOrg(org)
		r.SetName(name)
		r.SetFullName(fmt.Sprintf("%s/%s", org, name))

		// capture the count of repo secrets from the native service
		return c.Database.ListSecretsForRepo(r, nil, page, perPage)
	case constants.SecretShared:
		// check if we should capture secrets for multiple teams
		if name == "*" {
			c.Logger.WithFields(logrus.Fields{
				"org":   org,
				"teams": teams,
				"type":  sType,
			}).Tracef("listing native %s secrets for teams %s in org %s", sType, teams, org)

			// capture the count of shared secrets for multiple teams from the native service
			return c.Database.ListSecretsForTeams(org, teams, nil, page, perPage)
		}

		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"team": name,
			"type": sType,
		}).Tracef("listing native %s secrets for %s/%s", sType, org, name)

		// capture the count of shared secrets from the native service
		return c.Database.ListSecretsForTeam(org, name, nil, page, perPage)
	default:
		return nil, 0, fmt.Errorf("invalid secret type: %s", sType)
	}
}
