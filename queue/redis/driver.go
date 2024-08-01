// SPDX-License-Identifier: Apache-2.0

package redis

import "github.com/go-vela/types/constants"

// Driver outputs the configured queue driver.
func (c *client) Driver() string {
	return constants.DriverRedis
}
