// SPDX-License-Identifier: Apache-2.0

package vault

import "github.com/go-vela/server/constants"

// Driver outputs the configured secret driver.
func (c *client) Driver() string {
	return constants.DriverVault
}
