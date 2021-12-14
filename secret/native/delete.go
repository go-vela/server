// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

// Delete deletes a secret.
func (c *client) Delete(sType, org, name, path string) error {
	c.Logger.Tracef("deleting native %s secret %s for %s/%s", sType, path, org, name)

	// capture the secret from the native service
	s, err := c.Database.GetSecret(sType, org, name, path)
	if err != nil {
		return err
	}

	// delete the secret from the native service
	return c.Database.DeleteSecret(s.GetID())
}
