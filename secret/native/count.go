// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// Count counts a list of secrets.
func (c *Client) Count(ctx context.Context, sType, org, name string, teams []string) (int64, error) {
	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"type": sType,
		}).Tracef("counting native %s secrets for %s", sType, org)

		// capture the count of org secrets from the native service
		return c.Database.CountSecretsForOrg(ctx, org, nil)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"repo": name,
			"type": sType,
		}).Tracef("counting native %s secrets for %s/%s", sType, org, name)

		// create the repo with the information available
		r := new(api.Repo)
		r.SetOrg(org)
		r.SetName(name)
		r.SetFullName(fmt.Sprintf("%s/%s", org, name))

		// capture the count of repo secrets from the native service
		return c.Database.CountSecretsForRepo(ctx, r, nil)
	case constants.SecretShared:
		// check if we should capture secrets for multiple teams
		if name == "*" {
			c.Logger.WithFields(logrus.Fields{
				"org":   org,
				"teams": teams,
				"type":  sType,
			}).Tracef("counting native %s secrets for teams %s in org %s", sType, teams, org)

			// capture the count of shared secrets for multiple teams from the native service
			return c.Database.CountSecretsForTeams(ctx, org, teams, nil)
		}

		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"team": name,
			"type": sType,
		}).Tracef("counting native %s secrets for %s/%s", sType, org, name)

		// capture the count of shared secrets from the native service
		return c.Database.CountSecretsForTeam(ctx, org, name, nil)
	default:
		return 0, fmt.Errorf("invalid secret type: %s", sType)
	}
}
