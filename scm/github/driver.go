// SPDX-License-Identifier: Apache-2.0

package github

import "github.com/go-vela/types/constants"

// Driver outputs the configured scm driver.
func (c *client) Driver() string {
	return constants.DriverGithub
}
