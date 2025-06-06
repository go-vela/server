// SPDX-License-Identifier: Apache-2.0

package redis

import "github.com/go-vela/server/constants"

// Driver outputs the configured queue driver.
func (c *Client) Driver() string {
	return constants.DriverRedis
}
