// SPDX-License-Identifier: Apache-2.0

package native

import "github.com/go-vela/types/constants"

// Driver outputs the configured secret driver.
func (c *client) Driver() string {
	return constants.DriverNative
}
