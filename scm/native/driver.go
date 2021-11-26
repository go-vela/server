// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import "github.com/go-vela/types/constants"

// Driver outputs the configured scm driver.
func (c *client) Driver() string {
	// constants.DriverNative
	return constants.DriverNative
}
