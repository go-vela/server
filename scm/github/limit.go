// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// InstallRateLimit captures the SCM app rate limit for a given installation.
func (c *Client) InstallRateLimit(ctx context.Context, token string, installID int64) (int, int, int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"install_id": installID,
	}).Tracef("capturing SCM app rate limit for installation %d", installID)

	// create GitHub OAuth client with user's token
	client := c.newTokenClient(ctx, token)

	rateLimits, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("unable to get rate limits for installation %d: %w", installID, err)
	}

	return rateLimits.Core.Limit, rateLimits.Core.Remaining, rateLimits.Core.Reset.Time.Unix(), nil
}
