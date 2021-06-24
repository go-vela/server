// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package sqlite

import "github.com/go-vela/types/constants"

// Driver outputs the configured database driver.
func (c *client) Driver() string {
	return constants.DriverSqlite
}
