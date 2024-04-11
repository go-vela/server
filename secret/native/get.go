// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// Get captures a secret.
func (c *client) Get(ctx context.Context, sType, org, name, path string) (*library.Secret, error) {
	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": path,
			"type":   sType,
		}).Tracef("getting native %s secret %s for %s", sType, path, org)

		// capture the org secret from the native service
		return c.Database.GetSecretForOrg(ctx, org, path)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"repo":   name,
			"secret": path,
			"type":   sType,
		}).Tracef("getting native %s secret %s for %s/%s", sType, path, org, name)

		// create the repo with the information available
		r := new(api.Repo)
		r.SetOrg(org)
		r.SetName(name)
		r.SetFullName(fmt.Sprintf("%s/%s", org, name))

		// capture the repo secret from the native service
		return c.Database.GetSecretForRepo(ctx, path, r)
	case constants.SecretShared:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": path,
			"team":   name,
			"type":   sType,
		}).Tracef("getting native %s secret %s for %s/%s", sType, path, org, name)

		// capture the shared secret from the native service
		return c.Database.GetSecretForTeam(ctx, org, name, path)
	default:
		return nil, fmt.Errorf("invalid secret type: %s", sType)
	}
}
