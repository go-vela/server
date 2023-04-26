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

// Count counts a list of secrets.
func (c *client) Count(sType, org, name string, teams []string) (int64, error) {
	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"type": sType,
		}).Tracef("counting native %s secrets for %s", sType, org)

		// capture the count of org secrets from the native service
		return c.Database.CountSecretsForOrg(org, nil)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"repo": name,
			"type": sType,
		}).Tracef("counting native %s secrets for %s/%s", sType, org, name)

		// create the repo with the information available
		r := new(library.Repo)
		r.SetOrg(org)
		r.SetName(name)
		r.SetFullName(fmt.Sprintf("%s/%s", org, name))

		// capture the count of repo secrets from the native service
		return c.Database.CountSecretsForRepo(r, nil)
	case constants.SecretShared:
		// check if we should capture secrets for multiple teams
		if name == "*" {
			c.Logger.WithFields(logrus.Fields{
				"org":   org,
				"teams": teams,
				"type":  sType,
			}).Tracef("counting native %s secrets for teams %s in org %s", sType, teams, org)

			// capture the count of shared secrets for multiple teams from the native service
			return c.Database.CountSecretsForTeams(org, teams, nil)
		}

		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"team": name,
			"type": sType,
		}).Tracef("counting native %s secrets for %s/%s", sType, org, name)

		// capture the count of shared secrets from the native service
		return c.Database.CountSecretsForTeam(org, name, nil)
	default:
		return 0, fmt.Errorf("invalid secret type: %s", sType)
	}
}
