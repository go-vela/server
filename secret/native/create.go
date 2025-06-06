// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// Create creates a new secret.
func (c *Client) Create(ctx context.Context, sType, org, name string, s *api.Secret) (*api.Secret, error) {
	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("creating native %s secret %s for %s", sType, s.GetName(), org)

		// create the org secret in the native service
		return c.Database.CreateSecret(ctx, s)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"repo":   name,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("creating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

		// create the repo secret in the native service
		return c.Database.CreateSecret(ctx, s)
	case constants.SecretShared:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": s.GetName(),
			"team":   name,
			"type":   sType,
		}).Tracef("creating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

		// create the shared secret in the native service
		return c.Database.CreateSecret(ctx, s)
	default:
		return nil, fmt.Errorf("invalid secret type: %s", sType)
	}
}
