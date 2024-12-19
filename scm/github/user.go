// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// GetUserID captures the user's scm id.
func (c *client) GetUserID(ctx context.Context, name string, token string) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"user": name,
	}).Tracef("capturing SCM user id for %s", name)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, token)

	// send API call to capture user
	user, _, err := client.Users.Get(ctx, name)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(user.GetID()), nil
}
