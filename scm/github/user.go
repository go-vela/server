// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/sirupsen/logrus"
)

// GetUserID captures the user's scm id.
func (c *client) GetUserID(ctx context.Context, u *api.User) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("capturing scm id for %s", u.GetName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// send API call to capture user
	user, _, err := client.Users.Get(ctx, u.GetName())
	if err != nil {
		return "", err
	}

	return fmt.Sprint(user.GetID()), nil
}
