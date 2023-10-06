// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// List captures a list of secrets.
func (c *client) List(ctx context.Context, sType, org, name string, page, perPage int, teams []string) ([]*library.Secret, error) {
	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"type": sType,
		}).Tracef("listing native %s secrets for %s", sType, org)

		// capture the list of org secrets from the native service
		secrets, _, err := c.Database.ListSecretsForOrg(ctx, org, nil, page, perPage)
		if err != nil {
			return nil, err
		}

		return secrets, nil
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

		// capture the list of repo secrets from the native service
		secrets, _, err := c.Database.ListSecretsForRepo(ctx, r, nil, page, perPage)
		if err != nil {
			return nil, err
		}

		return secrets, nil
	case constants.SecretShared:
		// check if we should capture secrets for multiple teams
		if name == "*" {
			c.Logger.WithFields(logrus.Fields{
				"org":   org,
				"teams": teams,
				"type":  sType,
			}).Tracef("listing native %s secrets for teams %s in org %s", sType, teams, org)

			// capture the list of shared secrets for multiple teams from the native service
			secrets, _, err := c.Database.ListSecretsForTeams(ctx, org, teams, nil, page, perPage)
			if err != nil {
				return nil, err
			}

			return secrets, nil
		}

		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"team": name,
			"type": sType,
		}).Tracef("listing native %s secrets for %s/%s", sType, org, name)

		// capture the list of shared secrets from the native service
		secrets, _, err := c.Database.ListSecretsForTeam(ctx, org, name, nil, page, perPage)
		if err != nil {
			return nil, err
		}

		return secrets, nil
	default:
		return nil, fmt.Errorf("invalid secret type: %s", sType)
	}
}
