// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// GetOrgName gets org name from Github.
func (c *client) GetOrgName(ctx context.Context, u *api.User, o string) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  o,
		"user": u.GetName(),
	}).Tracef("retrieving org information for %s", o)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())

	// send an API call to get the org info
	orgInfo, resp, err := client.Organizations.Get(ctx, o)

	orgName := orgInfo.GetLogin()

	// if org is not found, return the personal org
	if resp.StatusCode == http.StatusNotFound {
		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			return "", err
		}

		orgName = user.GetLogin()
	} else if err != nil {
		return "", err
	}

	return orgName, nil
}
