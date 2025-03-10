// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// GetOrgIdentifiers gets org name and id from Github.
func (c *client) GetOrgIdentifiers(ctx context.Context, u *api.User, o string) (string, int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	}).Tracef("retrieving org information for %s", o)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())

	// send an API call to get the org info
	orgInfo, resp, err := client.Organizations.Get(ctx, o)

	orgName := orgInfo.GetLogin()
	orgID := orgInfo.GetID()

	// if org is not found, return the personal org
	if resp.StatusCode == http.StatusNotFound {
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			return "", 0, err
		}

		orgName = user.GetLogin()
		orgID = user.GetID()
	} else if err != nil {
		return "", 0, err
	}

	return orgName, orgID, nil
}
