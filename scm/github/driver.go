// SPDX-License-Identifier: Apache-2.0

package github

import "github.com/go-vela/server/constants"

// Driver outputs the configured scm driver.
func (c *Client) Driver() string {
	return constants.DriverGithub
}
