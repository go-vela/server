// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package vault

import "github.com/go-vela/types/constants"

// Driver outputs the configured secret driver.
func (c *client) Driver() string {
	return constants.DriverVault
}
