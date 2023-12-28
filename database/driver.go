// SPDX-License-Identifier: Apache-2.0

package database

// Driver outputs the configured database driver.
func (e *engine) Driver() string {
	return e.config.Driver
}
