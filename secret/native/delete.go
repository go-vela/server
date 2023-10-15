// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Delete deletes a secret.
func (c *client) Delete(ctx context.Context, sType, org, name, path string) error {
	// capture the secret from the native service
	s, err := c.Get(ctx, sType, org, name, path)
	if err != nil {
		return err
	}

	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": path,
			"type":   sType,
		}).Tracef("deleting native %s secret %s for %s", sType, path, org)

		// delete the org secret from the native service
		return c.Database.DeleteSecret(ctx, s)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"repo":   name,
			"secret": path,
			"type":   sType,
		}).Tracef("deleting native %s secret %s for %s/%s", sType, path, org, name)

		// delete the repo secret from the native service
		return c.Database.DeleteSecret(ctx, s)
	case constants.SecretShared:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": path,
			"team":   name,
			"type":   sType,
		}).Tracef("deleting native %s secret %s for %s/%s", sType, path, org, name)

		// delete the shared secret from the native service
		return c.Database.DeleteSecret(ctx, s)
	default:
		return fmt.Errorf("invalid secret type: %s", sType)
	}
}
