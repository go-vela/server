// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// Update updates an existing secret.
func (c *Client) Update(ctx context.Context, sType, org, name string, s *api.Secret) (*api.Secret, error) {
	// capture the secret from the native service
	secret, err := c.Get(ctx, sType, org, name, s.GetName())
	if err != nil {
		return nil, err
	}

	// update allow events if set
	if s.GetAllowEvents().ToDatabase() > 0 {
		secret.SetAllowEvents(s.GetAllowEvents())
	}

	// update the images if set
	if s.Images != nil {
		secret.SetImages(s.GetImages())
	}

	// update the value if set
	if len(s.GetValue()) > 0 {
		secret.SetValue(s.GetValue())
	}

	// update allow_command if set
	if s.AllowCommand != nil {
		secret.SetAllowCommand(s.GetAllowCommand())
	}

	// update allow_substitution if set
	if s.AllowSubstitution != nil {
		secret.SetAllowSubstitution(s.GetAllowSubstitution())
	}

	// update repo_allowlist if set
	if s.RepoAllowlist != nil {
		secret.SetRepoAllowlist(s.GetRepoAllowlist())
	}

	// update updated_at if set
	secret.SetUpdatedAt(s.GetUpdatedAt())

	// update updated_by if set
	secret.SetUpdatedBy(s.GetUpdatedBy())

	// handle the secret based off the type
	switch sType {
	case constants.SecretOrg:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("updating native %s secret %s for %s", sType, s.GetName(), org)

		// update the org secret in the native service
		return c.Database.UpdateSecret(ctx, secret)
	case constants.SecretRepo:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"repo":   name,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("updating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

		// update the repo secret in the native service
		return c.Database.UpdateSecret(ctx, secret)
	case constants.SecretShared:
		c.Logger.WithFields(logrus.Fields{
			"org":    org,
			"team":   name,
			"secret": s.GetName(),
			"type":   sType,
		}).Tracef("updating native %s secret %s for %s/%s", sType, s.GetName(), org, name)

		// update the shared secret in the native service
		return c.Database.UpdateSecret(ctx, secret)
	default:
		return nil, fmt.Errorf("invalid secret type: %s", sType)
	}
}
